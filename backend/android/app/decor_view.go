//go:build android

package app

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// decorView is the view the whole window is drawn into, which is what the
// system bars are asked about and told about.
func decorView(e *jni.Env) (jni.Object, error) {
	win, err := activityWindow(e)
	if err != nil {
		return jni.Object{}, err
	}
	winCls, err := e.Class("android/view/Window")
	if err != nil {
		return jni.Object{}, err
	}
	m, err := e.Method(winCls, "getDecorView", jni.Sig(jni.TClass("android/view/View")))
	if err != nil {
		return jni.Object{}, err
	}
	return e.CallObject(win, m)
}
