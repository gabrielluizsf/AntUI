//go:build android

package speech

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Speaking reports whether anything is being said now.
func (s *Speaker) Speaking() (bool, error) {
	var out bool
	err := jni.Do(func(e *jni.Env) error {
		var err error
		out, err = e.InvokeBool(s.tts, "isSpeaking", jni.Sig(jni.TBool))
		return err
	})
	return out, err
}
