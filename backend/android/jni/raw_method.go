//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"unsafe"
)

func (e *Env) rawMethod(cls C.jclass, name, sig string) C.jmethodID {
	cname, csig := C.CString(name), C.CString(sig)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(csig))
	m := C.mw_GetMethodID(e.ptr, cls, cname, csig)
	if m == nil {
		clearPending(e.ptr)
	}
	return m
}
