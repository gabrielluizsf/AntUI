//go:build android

package app

/*
#include "glue.h"
*/
import "C"

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// registerShim attaches the Go implementations to the shim's native methods.
//
// The obvious way is not to: a Go function exported as
// Java_dev_antui_AntuiActivity_nativeOnPermissionResult is the name the
// machine looks for, and it is in the library. It does not work, and the
// reason is worth writing down, because the error it produces —
// "UnsatisfiedLinkError: No implementation found" for a symbol that is
// plainly present in the .so — points nowhere near it.
//
// The machine only searches libraries that were loaded through
// System.loadLibrary, because that is what records which class loader a
// library belongs to. NativeActivity does not use it: it dlopens the library
// from its own native code. So the library is in the process, its symbols
// are exported, and as far as the machine is concerned it was never loaded
// at all. RegisterNatives says so explicitly, and is the only way.
func registerShim() {
	err := jni.Do(func(e *jni.Env) error {
		cls, err := e.Class(ShimClass)
		if err != nil {
			// No shim in this build. Not an error: an app that asks for
			// nothing from Java does not need one, and the calls that do
			// need it report ErrNoShim.
			return nil
		}
		return e.RegisterNatives(cls, shimNatives())
	})
	if err != nil {
		ndk.Errorf("the Java shim is in the package and could not be wired up: %v", err)
	}

	// The one door every listener in the shim comes back through; see
	// callback.go. It is on a class of its own, so it is a second call.
	if err := registerListenerBridge(); err != nil {
		ndk.Errorf("the listener bridge could not be wired up: %v", err)
	}
}
