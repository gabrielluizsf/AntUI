//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

//export antuiOnPause
func antuiOnPause(*C.ANativeActivity) { post(Pause, nil, true) }
