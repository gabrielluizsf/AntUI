//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// IsInstance reports whether the object is of that class, or of one
// descended from it.
func (e *Env) IsInstance(o Object, c Class) bool {
	return C.mw_IsInstanceOf(e.ptr, o.ref, c.ref) != C.JNI_FALSE
}
