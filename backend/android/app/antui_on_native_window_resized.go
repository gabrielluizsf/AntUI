//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

import "unsafe"

//export antuiOnNativeWindowResized
func antuiOnNativeWindowResized(_ *C.ANativeActivity, w *C.ANativeWindow) {
	post(WindowResized, unsafe.Pointer(w), false)
}
