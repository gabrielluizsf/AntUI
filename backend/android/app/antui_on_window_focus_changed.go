//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

//export antuiOnWindowFocusChanged
func antuiOnWindowFocusChanged(_ *C.ANativeActivity, hasFocus C.int) {
	if hasFocus != 0 {
		post(FocusGained, nil, false)
	} else {
		post(FocusLost, nil, false)
	}
}
