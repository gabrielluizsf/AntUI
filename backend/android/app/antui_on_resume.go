//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

//export antuiOnResume
func antuiOnResume(*C.ANativeActivity) { post(Resume, nil, true) }
