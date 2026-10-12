//go:build android

package intent

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Launch is the intent the app was started with, which is where a deep link
// arrives when the app was not already running. [app.OnNewIntent] is where
// it arrives when it was.
func Launch() (jni.Object, error) {
	var out jni.Object
	err := jni.Do(func(e *jni.Env) error {
		in, err := e.Invoke(app.Context(), "getIntent",
			jni.Sig(jni.TClass("android/content/Intent")))
		if err != nil {
			return err
		}
		if !in.IsNil() {
			out = e.Global(in)
		}
		return nil
	})
	return out, err
}
