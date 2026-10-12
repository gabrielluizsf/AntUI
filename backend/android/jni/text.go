//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// text calls a no-argument method that returns a String, and reads it. It
// swallows anything that goes wrong, for the same reason describe does.
func (e *Env) text(o C.jobject, name, sig string) string {
	cls := C.mw_GetObjectClass(e.ptr, o)
	if cls == 0 {
		clearPending(e.ptr)
		return ""
	}
	m := e.rawMethod(cls, name, sig)
	if m == nil {
		return ""
	}
	s := C.mw_CallObject(e.ptr, o, m, nil)
	clearPending(e.ptr)
	if s == 0 {
		return ""
	}
	out, err := e.GoString(Object{ref: s})
	if err != nil {
		return ""
	}
	return out
}
