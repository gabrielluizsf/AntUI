//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

func first(a []C.jvalue) *C.jvalue {
	if len(a) == 0 {
		return nil
	}
	return &a[0]
}
