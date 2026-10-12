//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// GetBool reads a boolean field of an object.
func (e *Env) GetBool(o Object, f Field) (bool, error) {
	if o.IsNil() {
		return false, errNilObject
	}
	r := C.mw_GetBooleanField(e.ptr, o.ref, f.id)
	return r != C.JNI_FALSE, e.check()
}
