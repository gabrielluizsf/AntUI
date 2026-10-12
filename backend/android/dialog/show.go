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

		text := func(method, value string) error {
			if value == "" {
				return nil
			}
			js, err := e.String(value)
			if err != nil {
				return err
			}
			_, err = e.Invoke(builder, method,
				jni.Sig(builderType, jni.TClass("java/lang/CharSequence")), jni.Ref(js))
			return err
		}
		if err := text("setTitle", title); err != nil {
			return err
		}
		if err := text("setMessage", message); err != nil {
			return err
		}
		for _, b := range []struct{ method, label string }{
			{"setPositiveButton", positive},
			{"setNegativeButton", negative},
			{"setNeutralButton", neutral},
		} {
			if b.label == "" {
				continue
			}
			js, err := e.String(b.label)
			if err != nil {
				return err
			}
			if _, err := e.Invoke(builder, b.method,
				jni.Sig(builderType, jni.TClass("java/lang/CharSequence"),
					jni.TClass("android/content/DialogInterface$OnClickListener")),
				jni.Ref(js), jni.Ref(listener)); err != nil {
				return err
			}
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
