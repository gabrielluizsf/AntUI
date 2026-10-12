//go:build android

package prefs

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Remove takes a key out.
func (p *Prefs) Remove(key string) error {
	return p.edit(func(e *jni.Env, ed jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		_, err = e.Invoke(ed, "remove",
			jni.Sig(jni.TClass("android/content/SharedPreferences$Editor"), jni.TString),
			jni.Ref(k))
		return err
	})
}
