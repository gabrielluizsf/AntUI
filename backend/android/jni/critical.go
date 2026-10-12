//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

// Critical hands f a view of a byte array's own storage, with no copy.
//
// It is worth it for a large array and only then. Between the two calls that
// bracket it the runtime may have stopped the collector, so **f must not
// call anything else in this package, must not block, and must not allocate
// anything that could make the collector run**. Breaking any of those is a
// deadlock inside the runtime rather than an error. When in doubt, use
// [Env.GoBytes] and pay for the copy.
func (e *Env) Critical(o Object, f func([]byte)) error {
	n, err := e.Len(o)
	if err != nil {
		return err
	}
	p := C.mw_GetCritical(e.ptr, C.jarray(o.ref))
	if p == nil {
		return errors.New("jni: cannot get at the array's storage")
	}
	f(unsafe.Slice((*byte)(p), n))
	// Mode 0: write anything back and free the copy, if the runtime made one.
	C.mw_ReleaseCritical(e.ptr, C.jarray(o.ref), p, 0)
	return e.check()
}
