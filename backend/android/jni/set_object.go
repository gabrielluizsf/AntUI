//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// SetObject writes an object field.
func (e *Env) SetObject(o Object, f Field, v Object) error {
	if o.IsNil() {
		return errNilObject
	}
	C.mw_SetObjectField(e.ptr, o.ref, f.id, v.ref)
	return e.check()
}
