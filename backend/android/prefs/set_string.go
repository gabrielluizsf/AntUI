//go:build android

package prefs

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// SetString, SetInt, SetBool, SetFloat and SetLong write one value.
func (p *Prefs) SetString(key, value string) error {
	return p.edit(func(e *jni.Env, ed jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		v, err := e.String(value)
		if err != nil {
			return err
		}
		_, err = e.Invoke(ed, "putString",
			jni.Sig(jni.TClass("android/content/SharedPreferences$Editor"),
				jni.TString, jni.TString), jni.Ref(k), jni.Ref(v))
		return err
	})
}
