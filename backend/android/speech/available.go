//go:build android

package speech

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Available reports whether anything on this device can listen.
func Available() bool {
	var out bool
	jni.Do(func(e *jni.Env) error {
		c, err := e.Class("android/speech/SpeechRecognizer")
		if err != nil {
			return err
		}
		m, err := e.StaticMethod(c, "isRecognitionAvailable",
			jni.Sig(jni.TBool, jni.TContext))
		if err != nil {
			return err
		}
		out, err = e.CallStaticBool(c, m, jni.Ref(app.Context()))
		return err
	})
	return out
}
