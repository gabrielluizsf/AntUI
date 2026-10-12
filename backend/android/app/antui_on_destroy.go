//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

//export antuiOnDestroy
func antuiOnDestroy(*C.ANativeActivity) {
	post(Destroy, nil, true)
	stateMu.Lock()
	theApp, activity = nil, nil
	stateMu.Unlock()
}
