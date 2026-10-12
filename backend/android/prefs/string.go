//go:build android

package prefs

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// String reads a string, or fallback when the key is not there.
func (p *Prefs) String(key, fallback string) (string, error) {
	var out string
	err := p.with(func(e *jni.Env, o jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		d, err := e.String(fallback)
		if err != nil {
			return err
		}
		out, err = e.InvokeString(o, "getString",
			jni.Sig(jni.TString, jni.TString, jni.TString), jni.Ref(k), jni.Ref(d))
		return err
	})
	return out, err
}
