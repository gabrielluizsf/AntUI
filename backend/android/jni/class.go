//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
	"strings"
	"sync"
	"unsafe"
)

// Class is a Java class, held as a global reference and cached, so that
// looking one up twice costs nothing the second time.
type Class struct{ ref C.jclass }

// IsNil reports whether the class was never found.
func (c Class) IsNil() bool { return c.ref == 0 }

// Object is the class itself as an object — a java.lang.Class instance,
// which is what a constructor like Intent(Context, Class) takes. A class
// reference and an object reference are the same thing to the machine; only
// Go's types tell them apart.
func (c Class) Object() Object { return Object{ref: C.jobject(c.ref)} }

// Method identifies a method of a class. It is not a reference and does not
// need freeing, but it is only valid while its class is loaded — which, for
// a cached class, is forever.
type Method struct{ id C.jmethodID }

// Field identifies a field, on the same terms as [Method].
type Field struct{ id C.jfieldID }

var (
	classMu sync.RWMutex
	classes = map[string]C.jclass{}
	// The app's own class loader, and its loadClass method. See Class.
	loader   C.jobject
	loadWith C.jmethodID
)

// Class finds a class by its name with slashes: "android/os/Build",
// "java/lang/String".
//
// It looks twice, and the second look is the point of this function.
// FindClass uses the class loader of whatever native code is on the stack —
// and on a thread this library attached, there is none, so it falls back to
// the *system* loader, which can see the framework and cannot see a single
// class from the app. The symptom is a ClassNotFoundException for a class
// that is plainly in the APK, on some threads and not others. The fix is to
// ask the app's own loader, which is what happens when the first look fails.
func (e *Env) Class(name string) (Class, error) {
	classMu.RLock()
	c, ok := classes[name]
	classMu.RUnlock()
	if ok {
		return Class{ref: c}, nil
	}

	ref := e.findClass(name)
	if ref == 0 {
		var err error
		if ref, err = e.loadClass(name); err != nil {
			return Class{}, err
		}
	}
	global := C.jclass(C.mw_NewGlobalRef(e.ptr, C.jobject(ref)))
	if global == 0 {
		return Class{}, errors.New("jni: cannot hold on to class " + name)
	}
	classMu.Lock()
	// Another goroutine may have got there first, in which case its
	// reference is the one everybody else already has.
	if was, ok := classes[name]; ok {
		classMu.Unlock()
		C.mw_DeleteGlobalRef(e.ptr, C.jobject(global))
		return Class{ref: was}, nil
	}
	classes[name] = global
	classMu.Unlock()
	return Class{ref: global}, nil
}

func (e *Env) findClass(name string) C.jclass {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ref := C.mw_FindClass(e.ptr, cname)
	if ref == 0 {
		// Not an error yet: the app's loader has not been tried.
		clearPending(e.ptr)
	}
	return ref
}

// loadClass asks the app's own class loader. Its argument is a dotted name,
// not a slashed one — the two spellings are used in different halves of the
// same API.
func (e *Env) loadClass(name string) (C.jclass, error) {
	l, m, err := e.appLoader()
	if err != nil {
		return 0, err
	}
	dotted, err := e.String(strings.ReplaceAll(name, "/", "."))
	if err != nil {
		return 0, err
	}
	a := jvalues([]Value{Ref(dotted)})
	got := C.mw_CallObject(e.ptr, l, m, first(a))
	keep(a)
	if err := e.check(); err != nil {
		return 0, errors.New("jni: no class " + name + ": " + err.Error())
	}
	if got == 0 {
		return 0, errors.New("jni: no class " + name)
	}
	return C.jclass(got), nil
}

// appLoader finds the class loader that loaded the activity, which is the
// one that can see the app's classes.
func (e *Env) appLoader() (C.jobject, C.jmethodID, error) {
	classMu.RLock()
	l, m := loader, loadWith
	classMu.RUnlock()
	if l != 0 {
		return l, m, nil
	}

	act := Activity()
	if act.IsNil() {
		return 0, nil, ErrNoVM
	}
	// The activity's class, then that class's own class, which is
	// java/lang/Class — reached this way rather than by name so that it
	// works on a thread FindClass would fail on.
	actClass := C.mw_GetObjectClass(e.ptr, act.ref)
	if actClass == 0 {
		return 0, nil, errors.New("jni: the activity has no class")
	}
	classClass := C.mw_GetObjectClass(e.ptr, C.jobject(actClass))
	get := e.rawMethod(classClass, "getClassLoader", "()Ljava/lang/ClassLoader;")
	if get == nil {
		return 0, nil, errors.New("jni: java.lang.Class has no getClassLoader")
	}
	ldr := C.mw_CallObject(e.ptr, C.jobject(actClass), get, nil)
	if err := e.check(); err != nil {
		return 0, nil, err
	}
	if ldr == 0 {
		return 0, nil, errors.New("jni: the activity's class has no loader")
	}
	ldrClass := C.mw_GetObjectClass(e.ptr, ldr)
	load := e.rawMethod(ldrClass, "loadClass", "(Ljava/lang/String;)Ljava/lang/Class;")
	if load == nil {
		return 0, nil, errors.New("jni: the class loader has no loadClass")
	}
	global := C.mw_NewGlobalRef(e.ptr, ldr)
	if global == 0 {
		return 0, nil, errors.New("jni: cannot hold on to the class loader")
	}

	classMu.Lock()
	if loader == 0 {
		loader, loadWith = global, load
	} else {
		C.mw_DeleteGlobalRef(e.ptr, global)
		global, load = loader, loadWith
	}
	l, m = loader, loadWith
	classMu.Unlock()
	return l, m, nil
}

// ClassOf is the class of an object, which is how a returned object's real
// type is found without knowing its name.
//
// **The caller owns what comes back and has to let it go**:
//
//	cls, err := e.ClassOf(o)
//	if err != nil {
//	    return err
//	}
//	defer e.Delete(cls.Object())
//
// It is a *local* reference, unlike the one [Env.Class] hands out, and two
// things follow from that. It has to be deleted, or a loop fills the local
// reference table — which is 512 entries deep and then the VM aborts. And it
// must never be used as the key of anything that outlives the call: the slot
// is reused, so the same number means one class now and a different class in
// a minute. Both of those were live bugs here, and neither showed up in a
// year of running until CheckJNI was turned on.
func (e *Env) ClassOf(o Object) (Class, error) {
	if o.IsNil() {
		return Class{}, errors.New("jni: no class for a null object")
	}
	ref := C.mw_GetObjectClass(e.ptr, o.ref)
	if ref == 0 {
		return Class{}, e.check()
	}
	return Class{ref: ref}, nil
}

// Method finds an instance method. The signature is built with [Sig].
func (e *Env) Method(c Class, name, sig string) (Method, error) {
	id, err := e.member(c, name, sig, func(cn, cs *C.char) unsafe.Pointer {
		return unsafe.Pointer(C.mw_GetMethodID(e.ptr, c.ref, cn, cs))
	}, "method")
	return Method{id: C.jmethodID(id)}, err
}

// StaticMethod finds a static method.
func (e *Env) StaticMethod(c Class, name, sig string) (Method, error) {
	id, err := e.member(c, name, sig, func(cn, cs *C.char) unsafe.Pointer {
		return unsafe.Pointer(C.mw_GetStaticMethodID(e.ptr, c.ref, cn, cs))
	}, "static method")
	return Method{id: C.jmethodID(id)}, err
}

// Constructor finds a constructor. Its name is "<init>" and it returns void,
// which is a rule rather than a convention, so the signature only needs its
// arguments: Constructor(cls, Sig(TVoid, TInt)).
func (e *Env) Constructor(c Class, sig string) (Method, error) {
	return e.Method(c, "<init>", sig)
}

// Field finds an instance field.
func (e *Env) Field(c Class, name, sig string) (Field, error) {
	id, err := e.member(c, name, sig, func(cn, cs *C.char) unsafe.Pointer {
		return unsafe.Pointer(C.mw_GetFieldID(e.ptr, c.ref, cn, cs))
	}, "field")
	return Field{id: C.jfieldID(id)}, err
}

// StaticField finds a static field, which is where most of the platform's
// constants live.
func (e *Env) StaticField(c Class, name, sig string) (Field, error) {
	id, err := e.member(c, name, sig, func(cn, cs *C.char) unsafe.Pointer {
		return unsafe.Pointer(C.mw_GetStaticFieldID(e.ptr, c.ref, cn, cs))
	}, "static field")
	return Field{id: C.jfieldID(id)}, err
}

// member is the half every lookup shares: the C strings, the nil check, and
// an error that says what was being looked for. A lookup that fails leaves
// a NoSuchMethodError pending, which has to be cleared or the next call
// fails instead.
func (e *Env) member(c Class, name, sig string, find func(cn, cs *C.char) unsafe.Pointer, what string) (unsafe.Pointer, error) {
	if c.IsNil() {
		return nil, errors.New("jni: no " + what + " " + name + ": the class was not found")
	}
	cname, csig := C.CString(name), C.CString(sig)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(csig))
	id := find(cname, csig)
	if id == nil {
		if err := e.check(); err != nil {
			return nil, errors.New("jni: no " + what + " " + name + sig + ": " + err.Error())
		}
		return nil, errors.New("jni: no " + what + " " + name + sig)
	}
	return id, nil
}
