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

// ShimClass is the Java class this library ships, if the app was built with
// it. An app built with -no-shim has none, and everything that depends on it
// reports [ErrNoShim].
const ShimClass = "dev/antui/AntuiActivity"

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
		return e.RegisterNatives(cls, []jni.Native{
			{
				Name: "nativeOnPermissionResult",
				Sig: jni.Sig(jni.TVoid, jni.TInt,
					jni.TArray(jni.TString), jni.TArray(jni.TInt)),
				Fn: C.antui_native_permission_result(),
			},
			{
				Name: "nativeOnActivityResult",
				Sig: jni.Sig(jni.TVoid, jni.TInt, jni.TInt,
					jni.TClass("android/content/Intent")),
				Fn: C.antui_native_activity_result(),
			},
			{
				Name: "nativeOnNewIntent",
				Sig:  jni.Sig(jni.TVoid, jni.TClass("android/content/Intent")),
				Fn:   C.antui_native_new_intent(),
			},
		})
	})
	if err != nil {
		ndk.Errorf("the Java shim is in the package and could not be wired up: %v", err)
	}

	// The one door every listener in the shim comes back through; see
	// callback.go. It is on a class of its own, so it is a second call.
	err = jni.Do(func(e *jni.Env) error {
		cls, err := e.Class("dev/antui/Native")
		if err != nil {
			return nil
		}
		return e.RegisterNatives(cls, []jni.Native{{
			Name: "invoke",
			Sig: jni.Sig(jni.TObject, jni.TLong, jni.TString,
				jni.TArray(jni.TObject)),
			Fn: C.antui_native_invoke(),
		}})
	})
	if err != nil {
		ndk.Errorf("the listener bridge could not be wired up: %v", err)
	}
}

// HasShim reports whether this app was built with the Java shim, and so
// whether permission results, activity results and arriving intents can be
// delivered at all.
func HasShim() bool {
	found := false
	jni.Do(func(e *jni.Env) error {
		_, err := e.Class(ShimClass)
		found = err == nil
		return nil
	})
	return found
}
