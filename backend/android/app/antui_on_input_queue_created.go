//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

import "unsafe"

//export antuiOnInputQueueCreated
func antuiOnInputQueueCreated(_ *C.ANativeActivity, q *C.AInputQueue) {
	post(InputUp, unsafe.Pointer(q), true)
}
