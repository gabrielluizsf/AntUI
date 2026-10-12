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

//export antuiOnNativeWindowRedrawNeeded
func antuiOnNativeWindowRedrawNeeded(_ *C.ANativeActivity, w *C.ANativeWindow) {
	// This one does not wait, and it is the one place that is a deliberate
	// loss. Waiting until the frame is actually on screen is what stops a
	// rotation flashing black — but the frame is drawn by the app between
	// pump and present, and holding the UI thread across that is holding it
	// across arbitrary user code, which may perfectly reasonably ask the UI
	// thread for the window insets and then wait for a thread that is
	// waiting for it. A flash is a blemish; a deadlock is the app gone.
	post(Redraw, unsafe.Pointer(w), false)
}
