//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
)

// Len is the length of a Java array.
func (e *Env) Len(o Object) (int, error) {
	if o.IsNil() {
		return 0, errors.New("jni: the array is null")
	}
	n := C.mw_GetArrayLength(e.ptr, C.jarray(o.ref))
	return int(n), e.check()
}
