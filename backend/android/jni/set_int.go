//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// SetInt writes an int field of an object.
func (e *Env) SetInt(o Object, f Field, v int32) error {
	if o.IsNil() {
		return errNilObject
	}
	C.mw_SetIntField(e.ptr, o.ref, f.id, C.jint(v))
	return e.check()
}
