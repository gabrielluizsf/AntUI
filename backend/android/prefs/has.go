//go:build android

package prefs

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Has reports whether a key is in the store at all, which is how a missing
// value is told from one that happens to equal the fallback.
func (p *Prefs) Has(key string) (bool, error) {
	var out bool
	err := p.with(func(e *jni.Env, o jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(o, "contains", jni.Sig(jni.TBool, jni.TString), jni.Ref(k))
		return err
	})
	return out, err
}
