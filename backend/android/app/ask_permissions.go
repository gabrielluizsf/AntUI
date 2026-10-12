//go:build android

package app

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// askPermissions calls the activity's requestPermissions, with the code the
// answer will carry.
func askPermissions(e *jni.Env, names []string, code int32) error {
	act := jni.Activity()
	cls, err := e.ClassOf(act)
	if err != nil {
		return err
	}
	defer e.Delete(cls.Object())
	m, err := e.Method(cls, "requestPermissions",
		jni.Sig(jni.TVoid, jni.TArray(jni.TString), jni.TInt))
	if err != nil {
		return err
	}
	arr, err := e.NewStrings(names)
	if err != nil {
		return err
	}
	return e.CallVoid(act, m, jni.Ref(arr), jni.Int(code))
}
