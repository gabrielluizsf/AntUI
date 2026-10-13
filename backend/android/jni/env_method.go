//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import "unsafe"

// Method finds an instance method. The signature is built with [Sig].
func (e *Env) Method(c Class, name, sig string) (Method, error) {
	id, err := e.member(c, name, sig, func(cn, cs *C.char) unsafe.Pointer {
		return unsafe.Pointer(C.mw_GetMethodID(e.ptr, c.ref, cn, cs))
	}, "method")
	return Method{id: C.jmethodID(id)}, err
}
