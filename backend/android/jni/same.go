//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// Same reports whether two references point at the same object. Two
// references to one object are not equal as values, so this is the only way
// to ask — and it is also how a weak reference is tested against nil.
func (e *Env) Same(a, b Object) bool {
	return C.mw_IsSameObject(e.ptr, a.ref, b.ref) != C.JNI_FALSE
}
