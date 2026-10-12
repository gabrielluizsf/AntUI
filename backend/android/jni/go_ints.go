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

// GoInts copies an int[] into Go.
func (e *Env) GoInts(o Object) ([]int32, error) {
	n, err := e.Len(o)
	if err != nil || n == 0 {
		return nil, err
	}
	out := make([]int32, n)
	C.mw_GetIntArrayRegion(e.ptr, C.jintArray(o.ref), 0, C.jsize(n), (*C.jint)(unsafe.Pointer(&out[0])))
	runtime.KeepAlive(out)
	return out, e.check()
}
