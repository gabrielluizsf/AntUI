//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// SetFloat writes a float field of an object.
func (e *Env) SetFloat(o Object, f Field, v float32) error {
	if o.IsNil() {
		return errNilObject
	}
	C.mw_SetFloatField(e.ptr, o.ref, f.id, C.jfloat(v))
	return e.check()
}
