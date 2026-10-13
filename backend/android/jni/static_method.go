//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import "unsafe"

// StaticMethod finds a static method.
func (e *Env) StaticMethod(c Class, name, sig string) (Method, error) {
	id, err := e.member(c, name, sig, func(cn, cs *C.char) unsafe.Pointer {
		return unsafe.Pointer(C.mw_GetStaticMethodID(e.ptr, c.ref, cn, cs))
	}, "static method")
	return Method{id: C.jmethodID(id)}, err
}
