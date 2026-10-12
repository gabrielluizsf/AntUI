//go:build android

package speech

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func (s *Speaker) tune(method string, v float32) error {
	return jni.Do(func(e *jni.Env) error {
		_, err := e.InvokeInt(s.tts, method, jni.Sig(jni.TInt, jni.TFloat), jni.Float(v))
		return err
	})
}
