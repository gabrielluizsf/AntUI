//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// GetFloat reads a float field of an object.
func (e *Env) GetFloat(o Object, f Field) (float32, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	r := C.mw_GetFloatField(e.ptr, o.ref, f.id)
	return float32(r), e.check()
}
