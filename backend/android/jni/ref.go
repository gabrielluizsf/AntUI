//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// Object is a Java object. It is a reference, not the object: which kind of
// reference decides how long it lasts.
//
//   - A **local** reference, which is what every call gives back, is valid
//     until the [Do] that made it returns. That is the whole point of Do:
//     the frame it opens is what frees them.
//   - A **global** reference, from [Env.Global], lasts until it is deleted,
//     and may be used from any thread. Everything cached between calls has
//     to be one.
//   - A **weak** reference does not keep the object alive, and may become
//     nil at any moment.
//
// Using a local reference after its frame has closed does not fail. It reads
// whatever now sits at that slot, which is another object, and the crash
// happens somewhere else entirely.
type Object struct {
	ref C.jobject
}

// Wrap takes a raw jobject from somewhere else — the activity handed over by
// the platform, an argument to a native method. The caller keeps whatever
// lifetime it already had.
//
// The handle is a uintptr because that is how cgo models Java's opaque
// types. It is not an address Go may follow, and nothing here ever follows
// one.
func Wrap(p uintptr) Object { return Object{ref: C.jobject(p)} }

// Handle is the raw reference, for handing to code that speaks JNI itself.
func (o Object) Handle() uintptr { return uintptr(o.ref) }

// IsNil reports whether there is no object — which is Java's null, and is
// what a method returns when it returns null.
func (o Object) IsNil() bool { return o.ref == 0 }

// Global makes a reference that outlives the call and may be used from any
// thread. Every one has to be given back with [Env.DeleteGlobal]; they do
// not go away on their own and the table they live in is not unlimited.
func (e *Env) Global(o Object) Object {
	if o.ref == 0 {
		return Object{}
	}
	return Object{ref: C.mw_NewGlobalRef(e.ptr, o.ref)}
}

// DeleteGlobal gives a global reference back.
func (e *Env) DeleteGlobal(o Object) {
	if o.ref != 0 {
		C.mw_DeleteGlobalRef(e.ptr, o.ref)
	}
}

// Weak makes a reference that does not keep the object alive. It may read as
// nil at any time, so anything using one has to check.
func (e *Env) Weak(o Object) Object {
	if o.ref == 0 {
		return Object{}
	}
	return Object{ref: C.mw_NewWeakGlobalRef(e.ptr, o.ref)}
}

// DeleteWeak gives a weak reference back.
func (e *Env) DeleteWeak(o Object) {
	if o.ref != 0 {
		C.mw_DeleteWeakGlobalRef(e.ptr, o.ref)
	}
}

// Delete frees one local reference early. It is only worth doing inside a
// loop that makes many — a frame frees the rest.
func (e *Env) Delete(o Object) {
	if o.ref != 0 {
		C.mw_DeleteLocalRef(e.ptr, o.ref)
	}
}

// Frame runs f inside a local reference frame of its own, so that everything
// f makes is freed when it returns. A loop that touches a thousand array
// elements needs this; anything else already has the frame [Do] opened.
func (e *Env) Frame(capacity int, f func() error) error {
	if C.mw_PushLocalFrame(e.ptr, C.jint(capacity)) != 0 {
		clearPending(e.ptr)
		return errNoFrame
	}
	err := f()
	C.mw_PopLocalFrame(e.ptr, 0)
	return err
}

// Same reports whether two references point at the same object. Two
// references to one object are not equal as values, so this is the only way
// to ask — and it is also how a weak reference is tested against nil.
func (e *Env) Same(a, b Object) bool {
	return C.mw_IsSameObject(e.ptr, a.ref, b.ref) != C.JNI_FALSE
}

// IsInstance reports whether the object is of that class, or of one
// descended from it.
func (e *Env) IsInstance(o Object, c Class) bool {
	return C.mw_IsInstanceOf(e.ptr, o.ref, c.ref) != C.JNI_FALSE
}
