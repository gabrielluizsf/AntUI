//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// ANativeActivity_onCreate is the one function Android looks for in this
// library. It runs on the UI thread and has to return quickly, so all it
// does is set the callbacks up and start the goroutine that is the app.
//
//export ANativeActivity_onCreate
func ANativeActivity_onCreate(act *C.ANativeActivity, savedState unsafe.Pointer, savedStateSize C.size_t) {
	// Before anything else, so that a panic in the lines below is readable
	// in logcat instead of vanishing.
	ndk.RedirectStdio()

	// This is the UI thread, and the only place its looper can be reached.
	startUIRunner()

	f := registeredMain()
	if f == nil {
		ndk.Errorf("no app registered: call app.Main(f) from an init function")
		return
	}

	a := buildApp(act, savedState, savedStateSize)

	stateMu.Lock()
	theApp, activity = a, act
	stateMu.Unlock()

	if err := initBridge(act, a); err != nil {
		ndk.Errorf("no bridge to Java: %v", err)
	}

	C.antui_set_callbacks(act)

	go func() {
		// The looper, the input queue and most of the NDK are thread-local,
		// so the app's goroutine keeps one thread for its whole life.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		run(a, f)
	}()
}
