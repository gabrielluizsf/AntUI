//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

//export antuiOnStart
func antuiOnStart(*C.ANativeActivity) { post(Start, nil, true) }
