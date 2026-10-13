//go:build android

package prefs

import "github.com/gabrielluizsf/antui/backend/android/jni"

// Float reads a float setting, or the fallback when it is not there.
func (p *Prefs) Float(key string, fallback float32) (float32, error) {
	var out float32
	err := p.with(func(e *jni.Env, o jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		out, err = e.InvokeFloat(o, "getFloat",
			jni.Sig(jni.TFloat, jni.TString, jni.TFloat), jni.Ref(k), jni.Float(fallback))
		return err
	})
	return out, err
}
