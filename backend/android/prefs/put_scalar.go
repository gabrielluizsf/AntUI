//go:build android

package prefs

import "github.com/gabrielluizsf/antui/backend/android/jni"

func (p *Prefs) putScalar(key, method, typ string, value jni.Value) error {
	return p.edit(func(e *jni.Env, ed jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		_, err = e.Invoke(ed, method,
			jni.Sig(jni.TClass("android/content/SharedPreferences$Editor"), jni.TString, typ),
			jni.Ref(k), value)
		return err
	})
}
