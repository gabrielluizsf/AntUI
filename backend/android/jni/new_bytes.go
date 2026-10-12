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

// NewBytes makes a byte[] holding a copy of v.
func (e *Env) NewBytes(v []byte) (Object, error) {
	a := C.mw_NewByteArray(e.ptr, C.jsize(len(v)))
	if err := e.check(); err != nil {
		return Object{}, err
	}
	if a == 0 {
		return Object{}, errors.New("jni: cannot make a byte array")
	}
	if len(v) > 0 {
		// A jbyte is signed and a Go byte is not, which is a difference in
		// how the number is read and not in the bits, so the pointer is
		// simply reinterpreted rather than the array copied twice.
		C.mw_SetByteArrayRegion(e.ptr, a, 0, C.jsize(len(v)), (*C.jbyte)(unsafe.Pointer(&v[0])))
		runtime.KeepAlive(v)
		if err := e.check(); err != nil {
			return Object{}, err
		}
	}
	return Object{ref: C.jobject(a)}, nil
}
