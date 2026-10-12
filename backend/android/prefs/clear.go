//go:build android

package prefs

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Clear empties the store.
func (p *Prefs) Clear() error {
	return p.edit(func(e *jni.Env, ed jni.Object) error {
		_, err := e.Invoke(ed, "clear",
			jni.Sig(jni.TClass("android/content/SharedPreferences$Editor")))
		return err
	})
}
