//go:build android

package speech

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// SetLanguage asks for a language by its tag — "pt-BR", "en-GB".
//
// It reports whether the engine has it. A language it knows but has not
// downloaded is as unusable as one it does not know, and both come back
// false: the difference is not something an app can do anything about.
func (s *Speaker) SetLanguage(tag string) (bool, error) {
	var ok bool
	err := jni.Do(func(e *jni.Env) error {
		jtag, err := e.String(tag)
		if err != nil {
			return err
		}
		locale, err := e.Static("java/util/Locale", "forLanguageTag",
			jni.Sig(jni.TClass("java/util/Locale"), jni.TString), jni.Ref(jtag))
		if err != nil {
			return err
		}
		r, err := e.InvokeInt(s.tts, "setLanguage",
			jni.Sig(jni.TInt, jni.TClass("java/util/Locale")), jni.Ref(locale))
		if err != nil {
			return err
		}
		ok = r != langMissingData && r != langNotSupported
		return nil
	})
	return ok, err
}
