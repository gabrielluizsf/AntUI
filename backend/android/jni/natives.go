//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

// Native is one native method to attach to a Java class.
//
// Registering by hand is rarely needed. A Go function exported as
// Java_pkg_Class_method is found by the machine on its own, by name, the
// first time the method is called — which is how the shim in
// antui/backend/android/app works and is one fewer thing to get wrong. This is for
// the cases where the name cannot be used: a class the build renamed, or a
// method whose Go implementation is chosen at run time.
type Native struct {
	// Name is the method's name in Java.
	Name string
	// Sig is its signature, from [Sig].
	Sig string
	// Fn is the C address of the implementation. It is not an ordinary Go
	// function: it has to be a symbol cgo exported, taken as C.theSymbol
	// from the package that exported it.
	Fn unsafe.Pointer
}

// RegisterNatives attaches implementations to a class's native methods.
//
// A signature that does not match one the class declares is not an error
// here — the machine reports it as a NoSuchMethodError, which is what comes
// back — but a signature that matches the *wrong* method is neither, and is
// a crash the first time it is called.
func (e *Env) RegisterNatives(c Class, methods []Native) error {
	if len(methods) == 0 {
		return nil
	}
	if c.IsNil() {
		return errors.New("jni: cannot register natives on a class that was not found")
	}
	size := C.size_t(unsafe.Sizeof(C.JNINativeMethod{}))
	raw := C.malloc(size * C.size_t(len(methods)))
	if raw == nil {
		return errors.New("jni: out of memory registering natives")
	}
	defer C.free(raw)

	table := unsafe.Slice((*C.JNINativeMethod)(raw), len(methods))
	var owned []unsafe.Pointer
	defer func() {
		for _, p := range owned {
			C.free(p)
		}
	}()
	for i, m := range methods {
		name, sig := C.CString(m.Name), C.CString(m.Sig)
		owned = append(owned, unsafe.Pointer(name), unsafe.Pointer(sig))
		table[i].name = name
		table[i].signature = sig
		table[i].fnPtr = m.Fn
	}

	if C.mw_RegisterNatives(e.ptr, c.ref, (*C.JNINativeMethod)(raw), C.jint(len(methods))) != C.JNI_OK {
		if err := e.check(); err != nil {
			return err
		}
		return errors.New("jni: registering natives failed")
	}
	return nil
}

// UnregisterNatives takes every native method off a class. The platform's
// own documentation calls it a tool for debugging and tool builders, and
// says an app should not need it; it is here for symmetry.
func (e *Env) UnregisterNatives(c Class) error {
	if C.mw_UnregisterNatives(e.ptr, c.ref) != C.JNI_OK {
		if err := e.check(); err != nil {
			return err
		}
		return errors.New("jni: unregistering natives failed")
	}
	return nil
}
