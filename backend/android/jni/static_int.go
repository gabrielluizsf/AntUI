//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// StaticInt reads an int field of a class rather than of an object.
func (e *Env) StaticInt(c Class, f Field) (int32, error) {
	r := C.mw_GetStaticIntField(e.ptr, c.ref, f.id)
	return int32(r), e.check()
}
