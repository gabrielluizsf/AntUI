//go:build android

package prefs

import "github.com/gabrielluizsf/antui/backend/android/jni"

// Bool reads a true-or-false setting, or the fallback when it is not
// there.
func (p *Prefs) Bool(key string, fallback bool) (bool, error) {
	var out bool
	err := p.with(func(e *jni.Env, o jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(o, "getBoolean",
			jni.Sig(jni.TBool, jni.TString, jni.TBool), jni.Ref(k), jni.Bool(fallback))
		return err
	})
	return out, err
}
