//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

//export antuiOnLowMemory
func antuiOnLowMemory(*C.ANativeActivity) { post(LowMemory, nil, false) }
