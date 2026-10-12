//go:build android

package speech

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Stop drops everything queued and stops mid-word.
func (s *Speaker) Stop() error {
	return jni.Do(func(e *jni.Env) error {
		_, err := e.InvokeInt(s.tts, "stop", jni.Sig(jni.TInt))
		return err
	})
}
