//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"unsafe"
)

// StaticField finds a static field, which is where most of the platform's
// constants live.
func (e *Env) StaticField(c Class, name, sig string) (Field, error) {
	id, err := e.member(c, name, sig, func(cn, cs *C.char) unsafe.Pointer {
		return unsafe.Pointer(C.mw_GetStaticFieldID(e.ptr, c.ref, cn, cs))
	}, "static field")
	return Field{id: C.jfieldID(id)}, err
}
