//go:build android

package app

import "github.com/gabrielluizsf/antui/backend/android/jni"

// FontScale is how much bigger the user has asked text to be: 1 is the
// default, and it goes past 2 on a device set up for poor eyesight. An app
// that draws its own text is the only thing that can honour it.
func FontScale() (float32, error) {
	var scale float32 = 1
	err := jni.Do(func(e *jni.Env) error {
		act := jni.Activity()
		actCls, err := e.ClassOf(act)
		if err != nil {
			return err
		}
		defer e.Delete(actCls.Object())
		getRes, err := e.Method(actCls, "getResources",
			jni.Sig(jni.TClass("android/content/res/Resources")))
		if err != nil {
			return err
		}
		res, err := e.CallObject(act, getRes)
		if err != nil {
			return err
		}
		resCls, err := e.Class("android/content/res/Resources")
		if err != nil {
			return err
		}
		getCfg, err := e.Method(resCls, "getConfiguration",
			jni.Sig(jni.TClass("android/content/res/Configuration")))
		if err != nil {
			return err
		}
		cfg, err := e.CallObject(res, getCfg)
		if err != nil {
			return err
		}
		cfgCls, err := e.Class("android/content/res/Configuration")
		if err != nil {
			return err
		}
		f, err := e.Field(cfgCls, "fontScale", jni.TFloat)
		if err != nil {
			return err
		}
		scale, err = e.GetFloat(cfg, f)
		return err
	})
	return scale, err
}
