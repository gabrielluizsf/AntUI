//go:build android

package app

/*
#include <android/native_activity.h>
*/
import "C"

import (
	"unsafe"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// initBridge wires the app to Java. This runs on the UI thread, which is the
// only place the activity reference handed over by the platform is valid —
// the first thing Init does is make it global, and everything else in the
// library reaches Java through that.
func initBridge(act *C.ANativeActivity, a *App) error {
	if err := jni.Init(unsafe.Pointer(act.vm), uintptr(act.clazz)); err != nil {
		return err
	}
	registerShim()
	// Before the app's goroutine starts, so that the first thing it does can
	// already make a temporary file.
	setupPaths(a)
	return nil
}
