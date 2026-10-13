//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import "unsafe"

// Field finds an instance field.
func (e *Env) Field(c Class, name, sig string) (Field, error) {
	id, err := e.member(c, name, sig, func(cn, cs *C.char) unsafe.Pointer {
		return unsafe.Pointer(C.mw_GetFieldID(e.ptr, c.ref, cn, cs))
	}, "field")
	return Field{id: C.jfieldID(id)}, err
}
