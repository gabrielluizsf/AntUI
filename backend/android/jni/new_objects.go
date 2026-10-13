//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import "errors"

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
