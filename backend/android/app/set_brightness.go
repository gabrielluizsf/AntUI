//go:build android

package app

import "github.com/gabrielluizsf/antui/backend/android/jni"

// SetBrightness overrides the screen's brightness while this app is in
// front: 0 for as dark as the device allows, 1 for full, and a negative
// number to hand the decision back to the system.
func SetBrightness(level float32) error {
	if level > 1 {
		level = 1
	}
	if level < 0 {
		level = -1 // the platform's "use the user's setting"
	}
	var err error
	if e := RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			win, err := activityWindow(e)
			if err != nil {
				return err
			}
			winCls, err := e.Class("android/view/Window")
			if err != nil {
				return err
			}
			paramsType := jni.TClass("android/view/WindowManager$LayoutParams")
			get, err := e.Method(winCls, "getAttributes", jni.Sig(paramsType))
			if err != nil {
				return err
			}
			params, err := e.CallObject(win, get)
			if err != nil {
				return err
			}
			pCls, err := e.Class("android/view/WindowManager$LayoutParams")
			if err != nil {
				return err
			}
			f, err := e.Field(pCls, "screenBrightness", jni.TFloat)
			if err != nil {
				return err
			}
			if err := e.SetFloat(params, f, level); err != nil {
				return err
			}
			set, err := e.Method(winCls, "setAttributes", jni.Sig(jni.TVoid, paramsType))
			if err != nil {
				return err
			}
			return e.CallVoid(win, set, jni.Ref(params))
		})
	}); e != nil {
		return e
	}
	return err
}
