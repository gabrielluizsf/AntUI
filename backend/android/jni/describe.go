//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// describe calls getClass on the throwable and then getName on the result.
// It uses the raw calls rather than the checked ones, because this is what
// runs when something has already gone wrong and a failure here must not
// turn into another error.
func (e *Env) describe(t C.jobject, get, sig, then string) string {
	cls := C.mw_GetObjectClass(e.ptr, t)
	if cls == 0 {
		clearPending(e.ptr)
		return ""
	}
	m := e.rawMethod(cls, get, sig)
	if m == nil {
		return ""
	}
	obj := C.mw_CallObject(e.ptr, t, m, nil)
	clearPending(e.ptr)
	if obj == 0 {
		return ""
	}
	return e.text(obj, then, "()Ljava/lang/String;")
}
