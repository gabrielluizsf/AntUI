//go:build android

package app

/*
#include "glue.h"
*/
import "C"

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// shimNatives is the activity's native methods and the Go functions that
// answer them.
func shimNatives() []jni.Native {
	return []jni.Native{
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
	}
}
