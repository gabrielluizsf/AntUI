//go:build android

package biometric

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// newPromptBuilder makes the prompt's builder with its title and subtitle.
func newPromptBuilder(e *jni.Env, title, subtitle string) (jni.Object, error) {
	builder, err := e.Make("android/hardware/biometrics/BiometricPrompt$Builder",
		jni.Sig(jni.TVoid, jni.TContext), jni.Ref(app.Context()))
	if err != nil {
		return jni.Object{}, err
	}
	bt := jni.TClass("android/hardware/biometrics/BiometricPrompt$Builder")
	cs := jni.TClass("java/lang/CharSequence")
	text := func(method, value string) error {
		if value == "" {
			return nil
		}
		js, err := e.String(value)
		if err != nil {
			return err
		}
		_, err = e.Invoke(builder, method, jni.Sig(bt, cs), jni.Ref(js))
		return err
	}
	if err := text("setTitle", title); err != nil {
		return jni.Object{}, err
	}
	if err := text("setSubtitle", subtitle); err != nil {
		return jni.Object{}, err
	}
	return builder, nil
}
