//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

import "unsafe"

//export antuiOnNativeWindowDestroyed
func antuiOnNativeWindowDestroyed(_ *C.ANativeActivity, w *C.ANativeWindow) {
	post(WindowDown, unsafe.Pointer(w), true)
}
