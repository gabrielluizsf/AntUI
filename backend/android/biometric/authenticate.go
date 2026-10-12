//go:build android

package biometric

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func authenticate(token int64, title, subtitle, cancel string) error {
	return jni.Do(func(e *jni.Env) error {
		builder, err := e.Make("android/hardware/biometrics/BiometricPrompt$Builder",
			jni.Sig(jni.TVoid, jni.TContext), jni.Ref(app.Context()))
		if err != nil {
			return err
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
			return err
		}
		if err := text("setSubtitle", subtitle); err != nil {
			return err
		}

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
