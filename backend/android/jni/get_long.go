//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// GetLong reads a long field of an object.
func (e *Env) GetLong(o Object, f Field) (int64, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	r := C.mw_GetLongField(e.ptr, o.ref, f.id)
	return int64(r), e.check()
}
