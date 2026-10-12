//go:build android

package intent

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Data is the address an intent carries — the link that was followed, the
// file that was opened. It is empty when there is none.
func Data(in jni.Object) (string, error) {
	if in.IsNil() {
		return "", nil
	}
	var out string
	err := jni.Do(func(e *jni.Env) error {
		var err error
		out, err = e.InvokeString(in, "getDataString", jni.Sig(jni.TString))
		return err
	})
	return out, err
}
