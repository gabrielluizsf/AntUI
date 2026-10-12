//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"unsafe"
)

// jvalues lays the arguments out the way the machine expects to read them.
func jvalues(vals []Value) []C.jvalue {
	if len(vals) == 0 {
		return nil
	}
	out := make([]C.jvalue, len(vals))
	for i, v := range vals {
		*(*uint64)(unsafe.Pointer(&out[i])) = uint64(v)
	}
	return out
}
