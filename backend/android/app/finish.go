//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

// Finish asks Android to close the activity. It returns immediately; the
// Destroy event arrives later, in the usual order after Pause and Stop.
func Finish() {
	stateMu.Lock()
	act := activity
	stateMu.Unlock()
	if act != nil {
		C.ANativeActivity_finish(act)
	}
}
