//go:build android

package app

import "github.com/gabrielluizsf/antui/backend/android/jni"

func setSystemUI(flags int) error {
	var err error
	if e := RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			decor, err := decorView(e)
			if err != nil {
				return err
			}
			viewCls, err := e.Class("android/view/View")
			if err != nil {
				return err
			}
			m, err := e.Method(viewCls, "setSystemUiVisibility", jni.Sig(jni.TVoid, jni.TInt))
			if err != nil {
				return err
			}
			return e.CallVoid(decor, m, jni.Int(int32(flags)))
		})
	}); e != nil {
		return e
	}
	return err
}
