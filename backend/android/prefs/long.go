//go:build android

package prefs

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Long reads a 64-bit setting, or the fallback when it is not there.
func (p *Prefs) Long(key string, fallback int64) (int64, error) {
	var out int64
	err := p.with(func(e *jni.Env, o jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		out, err = e.InvokeLong(o, "getLong",
			jni.Sig(jni.TLong, jni.TString, jni.TLong), jni.Ref(k), jni.Long(fallback))
		return err
	})
	return out, err
}
