//go:build android

package app

/*
#include <android/native_activity.h>
*/
import "C"

import "unsafe"

// buildApp is the state the app is created with.
func buildApp(act *C.ANativeActivity, savedState unsafe.Pointer, savedStateSize C.size_t) *App {
	return &App{
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
}
