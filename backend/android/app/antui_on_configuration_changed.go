//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

//export antuiOnConfigurationChanged
func antuiOnConfigurationChanged(act *C.ANativeActivity) {
	// A configuration change is where the density moves: a window dragged to
	// another display, or a device that reports differently once it is
	// unfolded.
	if a := current(); a != nil {
		cfg := readConfig(act)
		a.mu.Lock()
		a.config = cfg
		a.mu.Unlock()
	}
	post(ConfigChanged, nil, false)
}
