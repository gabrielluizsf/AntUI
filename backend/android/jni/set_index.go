//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import "errors"

// SetIndex writes one element of an object array.
func (e *Env) SetIndex(o Object, i int, v Object) error {
	if o.IsNil() {
		return errors.New("jni: the array is null")
	}
	C.mw_SetObjectArrayElement(e.ptr, C.jobjectArray(o.ref), C.jsize(i), v.ref)
	return e.check()
}
