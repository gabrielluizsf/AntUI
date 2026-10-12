//go:build android

package intent

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Action is what an intent asks for: [ActionView] for a link that was
// opened, and so on.
func Action(in jni.Object) (string, error) {
	if in.IsNil() {
		return "", nil
	}
	var out string
	err := jni.Do(func(e *jni.Env) error {
		var err error
		out, err = e.InvokeString(in, "getAction", jni.Sig(jni.TString))
		return err
	})
	return out, err
}
