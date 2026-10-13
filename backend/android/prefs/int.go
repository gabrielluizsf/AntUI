//go:build android

package prefs

import "github.com/gabrielluizsf/antui/backend/android/jni"

// Int, Bool, Float and Long are String for the other types the store holds.
func (p *Prefs) Int(key string, fallback int32) (int32, error) {
	var out int32
	err := p.with(func(e *jni.Env, o jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		out, err = e.InvokeInt(o, "getInt",
			jni.Sig(jni.TInt, jni.TString, jni.TInt), jni.Ref(k), jni.Int(fallback))
		return err
	})
	return out, err
}
