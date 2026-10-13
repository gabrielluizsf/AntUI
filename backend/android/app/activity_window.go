//go:build android

package app

import "github.com/gabrielluizsf/antui/backend/android/jni"

// activityWindow is the activity's android.view.Window.
func activityWindow(e *jni.Env) (jni.Object, error) {
	act := jni.Activity()
	cls, err := e.ClassOf(act)
	if err != nil {
		return jni.Object{}, err
	}
	defer e.Delete(cls.Object())
	m, err := e.Method(cls, "getWindow", jni.Sig(jni.TClass("android/view/Window")))
	if err != nil {
		return jni.Object{}, err
	}
	return e.CallObject(act, m)
}
