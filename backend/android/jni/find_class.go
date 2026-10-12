//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"unsafe"
)

func (e *Env) findClass(name string) C.jclass {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ref := C.mw_FindClass(e.ptr, cname)
	if ref == 0 {
		// Not an error yet: the app's loader has not been tried.
		clearPending(e.ptr)
	}
	return ref
}
