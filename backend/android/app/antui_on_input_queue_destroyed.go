//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

import (
	"unsafe"
)

//export antuiOnInputQueueDestroyed
func antuiOnInputQueueDestroyed(_ *C.ANativeActivity, q *C.AInputQueue) {
	post(InputDown, unsafe.Pointer(q), true)
}
