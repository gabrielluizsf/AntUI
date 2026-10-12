//go:build android

package app

/*
#cgo LDFLAGS: -landroid

#include <android/configuration.h>
#include <android/native_activity.h>
*/
import "C"

// SetWindowFlags adds and removes flags on the activity's window. It may be
// called from any thread; the platform posts it to the UI thread itself.
func SetWindowFlags(add, remove uint32) {
	stateMu.Lock()
	act := activity
	stateMu.Unlock()
	if act != nil {
		C.ANativeActivity_setWindowFlags(act, C.uint32_t(add), C.uint32_t(remove))
	}
}
