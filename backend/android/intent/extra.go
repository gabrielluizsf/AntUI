//go:build android

package intent

import "github.com/gabrielluizsf/antui/backend/android/jni"

// Extra reads one string an intent was given — which is how a notification
// says which notification it was.
func Extra(in jni.Object, key string) (string, error) {
	if in.IsNil() {
		return "", nil
	}
	var out string
	err := jni.Do(func(e *jni.Env) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		out, err = e.InvokeString(in, "getStringExtra",
			jni.Sig(jni.TString, jni.TString), jni.Ref(k))
		return err
	})
	return out, err
}
