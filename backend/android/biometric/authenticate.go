//go:build android

package biometric

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func authenticate(token int64, title, subtitle, cancel string) error {
	return jni.Do(func(e *jni.Env) error {
		builder, err := newPromptBuilder(e, title, subtitle)
		if err != nil {
			return err
		}
		bt := jni.TClass("android/hardware/biometrics/BiometricPrompt$Builder")
		cs := jni.TClass("java/lang/CharSequence")

		// The executor the platform runs its callbacks on. The main one will
		// do: everything it calls hands straight over to a channel.
		executor, err := e.Invoke(app.Context(), "getMainExecutor",
			jni.Sig(jni.TClass("java/util/concurrent/Executor")))
		if err != nil {
			return err
		}
		listener, err := app.Proxy(e, token,
			"android.content.DialogInterface$OnClickListener")
		if err != nil {
			return err
		}
		jcancel, err := e.String(cancel)
		if err != nil {
			return err
		}
		if _, err := e.Invoke(builder, "setNegativeButton",
			jni.Sig(bt, cs, jni.TClass("java/util/concurrent/Executor"),
				jni.TClass("android/content/DialogInterface$OnClickListener")),
			jni.Ref(jcancel), jni.Ref(executor), jni.Ref(listener)); err != nil {
			return err
		}
		prompt, err := e.Invoke(builder, "build",
			jni.Sig(jni.TClass("android/hardware/biometrics/BiometricPrompt")))
		if err != nil {
			return err
		}
		signal, err := e.Make("android/os/CancellationSignal", jni.Sig(jni.TVoid))
		if err != nil {
			return err
		}
		callback, err := app.Listener(e, app.AuthCallbackClass, token)
		if err != nil {
			return err
		}
		return e.InvokeVoid(prompt, "authenticate",
			jni.Sig(jni.TVoid, jni.TClass("android/os/CancellationSignal"),
				jni.TClass("java/util/concurrent/Executor"),
				jni.TClass("android/hardware/biometrics/BiometricPrompt$AuthenticationCallback")),
			jni.Ref(signal), jni.Ref(executor), jni.Ref(callback))
	})
}
