//go:build android

package dialog

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func show(token int64, title, message, positive, negative, neutral string) error {
	return jni.Do(func(e *jni.Env) error {
		builder, err := e.Make("android/app/AlertDialog$Builder",
			jni.Sig(jni.TVoid, jni.TContext), jni.Ref(app.Context()))
		if err != nil {
			return err
		}
		listener, err := app.Proxy(e, token, listenerInterfaces...)
		if err != nil {
			return err
		}
		if err := dialogText(e, builder, "setTitle", title); err != nil {
			return err
		}
		if err := dialogText(e, builder, "setMessage", message); err != nil {
			return err
		}
		if err := dialogButtons(e, builder, listener, positive, negative, neutral); err != nil {
			return err
		}
		// The back button and a tap outside both have to answer, or a caller
		// waiting on the channel waits forever.
		if _, err := e.Invoke(builder, "setOnCancelListener",
			jni.Sig(builderType,
				jni.TClass("android/content/DialogInterface$OnCancelListener")),
			jni.Ref(listener)); err != nil {
			return err
		}
		_, err = e.Invoke(builder, "show",
			jni.Sig(jni.TClass("android/app/AlertDialog")))
		return err
	})
}
