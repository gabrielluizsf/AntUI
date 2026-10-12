//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

// NewInts makes an int[] holding a copy of v.
func (e *Env) NewInts(v []int32) (Object, error) {
	a := C.mw_NewIntArray(e.ptr, C.jsize(len(v)))
	if err := e.check(); err != nil {
		return Object{}, err
	}
	if a == 0 {
		return Object{}, errors.New("jni: cannot make an int array")
	}
	if len(v) > 0 {
		C.mw_SetIntArrayRegion(e.ptr, a, 0, C.jsize(len(v)), (*C.jint)(unsafe.Pointer(&v[0])))
		runtime.KeepAlive(v)
		if err := e.check(); err != nil {
			return Object{}, err
		}
	}
	return Object{ref: C.jobject(a)}, nil
}
