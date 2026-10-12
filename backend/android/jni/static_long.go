//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// StaticLong reads a long field of a class.
func (e *Env) StaticLong(c Class, f Field) (int64, error) {
	r := C.mw_GetStaticLongField(e.ptr, c.ref, f.id)
	return int64(r), e.check()
}
