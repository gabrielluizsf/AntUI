//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// GoBytes copies a byte[] into Go.
func (e *Env) GoBytes(o Object) ([]byte, error) {
	n, err := e.Len(o)
	if err != nil || n == 0 {
		return nil, err
	}
	out := make([]byte, n)
	C.mw_GetByteArrayRegion(e.ptr, C.jbyteArray(o.ref), 0, C.jsize(n), (*C.jbyte)(unsafe.Pointer(&out[0])))
	runtime.KeepAlive(out)
	return out, e.check()
}
