//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

import "unsafe"

//export antuiOnNativeWindowCreated
func antuiOnNativeWindowCreated(_ *C.ANativeActivity, w *C.ANativeWindow) {
	post(WindowUp, unsafe.Pointer(w), true)
}
