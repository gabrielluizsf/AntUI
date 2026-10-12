//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
)

// Index reads one element of an object array.
func (e *Env) Index(o Object, i int) (Object, error) {
	if o.IsNil() {
		return Object{}, errors.New("jni: the array is null")
	}
	r := C.mw_GetObjectArrayElement(e.ptr, C.jobjectArray(o.ref), C.jsize(i))
	return Object{ref: r}, e.check()
}
