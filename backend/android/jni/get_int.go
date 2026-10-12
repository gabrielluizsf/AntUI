//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// GetInt reads an int field of an object.
func (e *Env) GetInt(o Object, f Field) (int32, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	r := C.mw_GetIntField(e.ptr, o.ref, f.id)
	return int32(r), e.check()
}
