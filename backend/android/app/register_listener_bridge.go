//go:build android

package app

/*
#include "glue.h"
*/
import "C"

import "github.com/gabrielluizsf/antui/backend/android/jni"

// registerListenerBridge wires the one method every listener comes back
// through, on the class of its own.
func registerListenerBridge() error {
	return jni.Do(func(e *jni.Env) error {
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
}
