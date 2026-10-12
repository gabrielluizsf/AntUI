//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// StaticObject reads an object field of a class, as a local reference.
func (e *Env) StaticObject(c Class, f Field) (Object, error) {
	r := C.mw_GetStaticObjectField(e.ptr, c.ref, f.id)
	return Object{ref: r}, e.check()
}
