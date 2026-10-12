//go:build android

package app

/*
#cgo LDFLAGS: -landroid

#include <android/configuration.h>
#include <android/native_activity.h>
*/
import "C"

// HideKeyboard puts the on-screen keyboard away.
func HideKeyboard(flags uint32) {
	stateMu.Lock()
	act := activity
	stateMu.Unlock()
	if act != nil {
		C.ANativeActivity_hideSoftInput(act, C.uint32_t(flags))
	}
}
