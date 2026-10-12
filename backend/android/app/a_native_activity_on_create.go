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

	"github.com/gabrielluizsf/antui/backend/android/jni"
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

	a := &App{
		// A little buffering, so a burst of configuration changes does not
		// block the UI thread on an app that is mid-frame.
		msgs: make(chan message, 16),
		dead: make(chan struct{}),
		info: Info{
			InternalData: C.GoString(act.internalDataPath),
			ExternalData: C.GoString(act.externalDataPath),
			OBB:          C.GoString(act.obbPath),
			SDK:          int(act.sdkVersion),
		},
		config: readConfig(act),
		assets: unsafe.Pointer(act.assetManager),
		state:  Created,
		saved:  copyState(savedState, savedStateSize),
	}

	stateMu.Lock()
	theApp, activity = a, act
	stateMu.Unlock()

	// The bridge to Java, before anything else can want it. This runs on the
	// UI thread, which is the only place the activity reference handed over
	// by the platform is valid — the first thing Init does is make it
	// global, and everything else in the library reaches Java through that.
	if err := jni.Init(unsafe.Pointer(act.vm), uintptr(act.clazz)); err != nil {
		ndk.Errorf("no bridge to Java: %v", err)
	} else {
		registerShim()
		// Before the app's goroutine starts, so that the first thing it does
		// can already make a temporary file.
		setupPaths(a)
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
