//go:build android

package jni

/*
#cgo LDFLAGS: -landroid -llog

#include "jni_wrap.h"
*/
import "C"

// clearPending drops an exception without looking at it. It is only for the
// places where looking would need another call and that call is exactly what
// cannot be made.
func clearPending(env *C.JNIEnv) {
	if C.mw_ExceptionCheck(env) != C.JNI_FALSE {
		C.mw_ExceptionClear(env)
	}
}
