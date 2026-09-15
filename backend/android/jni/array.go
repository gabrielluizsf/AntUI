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

// Len is the length of a Java array.
func (e *Env) Len(o Object) (int, error) {
	if o.IsNil() {
		return 0, errors.New("jni: the array is null")
	}
	n := C.mw_GetArrayLength(e.ptr, C.jarray(o.ref))
	return int(n), e.check()
}

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

// NewObjects makes an array of n objects of class c, all null.
func (e *Env) NewObjects(c Class, n int) (Object, error) {
	a := C.mw_NewObjectArray(e.ptr, C.jsize(n), c.ref, 0)
	if err := e.check(); err != nil {
		return Object{}, err
	}
	if a == 0 {
		return Object{}, errors.New("jni: cannot make an object array")
	}
	return Object{ref: C.jobject(a)}, nil
}

// Index reads one element of an object array.
func (e *Env) Index(o Object, i int) (Object, error) {
	if o.IsNil() {
		return Object{}, errors.New("jni: the array is null")
	}
	r := C.mw_GetObjectArrayElement(e.ptr, C.jobjectArray(o.ref), C.jsize(i))
	return Object{ref: r}, e.check()
}

// SetIndex writes one element of an object array.
func (e *Env) SetIndex(o Object, i int, v Object) error {
	if o.IsNil() {
		return errors.New("jni: the array is null")
	}
	C.mw_SetObjectArrayElement(e.ptr, C.jobjectArray(o.ref), C.jsize(i), v.ref)
	return e.check()
}

// NewStrings makes a String[] out of Go strings.
func (e *Env) NewStrings(v []string) (Object, error) {
	cls, err := e.Class("java/lang/String")
	if err != nil {
		return Object{}, err
	}
	arr, err := e.NewObjects(cls, len(v))
	if err != nil {
		return Object{}, err
	}
	// A frame of its own: each string is a local reference, and a list of a
	// few hundred would otherwise fill the table the caller is sharing.
	err = e.Frame(len(v)+8, func() error {
		for i, s := range v {
			js, err := e.String(s)
			if err != nil {
				return err
			}
			if err := e.SetIndex(arr, i, js); err != nil {
				return err
			}
		}
		return nil
	})
	return arr, err
}

// GoStrings reads a String[] into Go.
func (e *Env) GoStrings(o Object) ([]string, error) {
	n, err := e.Len(o)
	if err != nil || n == 0 {
		return nil, err
	}
	out := make([]string, n)
	err = e.Frame(n+8, func() error {
		for i := range n {
			el, err := e.Index(o, i)
			if err != nil {
				return err
			}
			if out[i], err = e.GoString(el); err != nil {
				return err
			}
		}
		return nil
	})
	return out, err
}

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
