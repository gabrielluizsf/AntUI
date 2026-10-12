//go:build android

package app

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func readInsets() (Insets, error) {
	var out Insets
	err := jni.Do(func(e *jni.Env) error {
		decor, err := decorView(e)
		if err != nil {
			return err
		}
		viewCls, err := e.Class("android/view/View")
		if err != nil {
			return err
		}
		// getRootWindowInsets arrived in API 23. Below that there is no way
		// to ask at all, so the honest answer is nothing rather than a guess.
		get, err := e.Method(viewCls, "getRootWindowInsets",
			jni.Sig(jni.TClass("android/view/WindowInsets")))
		if err != nil {
			return nil
		}
		ins, err := e.CallObject(decor, get)
		if err != nil || ins.IsNil() {
			return err
		}
		wiCls, err := e.Class("android/view/WindowInsets")
		if err != nil {
			return err
		}
		for _, side := range []struct {
			name string
			into *int
		}{
			{"getSystemWindowInsetLeft", &out.Left},
			{"getSystemWindowInsetTop", &out.Top},
			{"getSystemWindowInsetRight", &out.Right},
			{"getSystemWindowInsetBottom", &out.Bottom},
		} {
			m, err := e.Method(wiCls, side.name, jni.Sig(jni.TInt))
			if err != nil {
				return err
			}
			v, err := e.CallInt(ins, m)
			if err != nil {
				return err
			}
			*side.into = int(v)
		}
		out.Keyboard = imeInset(e, ins, wiCls)
		return nil
	})
	return out, err
}
