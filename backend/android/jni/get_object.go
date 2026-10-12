//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// GetObject reads an object field, as a local reference.
func (e *Env) GetObject(o Object, f Field) (Object, error) {
	if o.IsNil() {
		return Object{}, errNilObject
	}
	r := C.mw_GetObjectField(e.ptr, o.ref, f.id)
	return Object{ref: r}, e.check()
}
