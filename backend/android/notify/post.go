//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Post puts a notification in the shade, or replaces the one with the same
// id.
//
// From API 33 this needs the POST_NOTIFICATIONS permission, and without it
// nothing appears and nothing fails. Ask for it with
// [antui/backend/android/permission.Ensure] before the first post.
func Post(n Notification) error {
	a := app.Current()
	if a == nil {
		return app.ErrNoUIThread
	}
	sdk := a.Config().SDK
	return jni.Do(func(e *jni.Env) error {
		builder, err := newBuilder(e, sdk, n.Channel)
		if err != nil {
			return err
		}
		set := func(method, sig string, args ...jni.Value) error {
			_, err := e.Invoke(builder, method, sig, args...)
			return err
		}
		title, err := e.String(n.Title)
		if err != nil {
			return err
		}
		text, err := e.String(n.Text)
		if err != nil {
			return err
		}
		cs := jni.TClass("java/lang/CharSequence")
		b := jni.TClass("android/app/Notification$Builder")
		if err := set("setContentTitle", jni.Sig(b, cs), jni.Ref(title)); err != nil {
			return err
		}
		if err := set("setContentText", jni.Sig(b, cs), jni.Ref(text)); err != nil {
			return err
		}

		icon := n.Icon
		if icon == 0 {
			// Any notification without a small icon is refused outright, so
			// there has to be a default, and the only one every device has
			// is the platform's own.
			if icon, err = e.ConstantInt("android/R$drawable", "ic_dialog_info"); err != nil {
				return err
			}
		}
		if err := set("setSmallIcon", jni.Sig(b, jni.TInt), jni.Int(icon)); err != nil {
			return err
		}
		if n.AutoCancel {
			if err := set("setAutoCancel", jni.Sig(b, jni.TBool), jni.Bool(true)); err != nil {
				return err
			}
		}
		if n.Ongoing {
			if err := set("setOngoing", jni.Sig(b, jni.TBool), jni.Bool(true)); err != nil {
				return err
			}
		}
		if n.Total > 0 || n.Indeterminate {
			if err := set("setProgress", jni.Sig(b, jni.TInt, jni.TInt, jni.TBool),
				jni.Int(n.Total), jni.Int(n.Progress), jni.Bool(n.Indeterminate)); err != nil {
				return err
			}
		}
		if n.Silent && sdk >= 26 {
			// Below 26 the sound is a property of the notification and above
			// it a property of the channel, so this only means anything one
			// side of that line.
			if err := set("setSilent", jni.Sig(b, jni.TBool), jni.Bool(true)); err != nil {
				// setSilent is API 29. Not having it is not a failure.
				_ = err
			}
		}

		if err := addActions(e, builder, n); err != nil {
			return err
		}

		pending, err := tapIntent(e, n.ID, n.Extra)
		if err != nil {
			return err
		}
		if !pending.IsNil() {
			if err := set("setContentIntent",
				jni.Sig(b, jni.TClass("android/app/PendingIntent")), jni.Ref(pending)); err != nil {
				return err
			}
		}

		built, err := e.Invoke(builder, "build",
			jni.Sig(jni.TClass("android/app/Notification")))
		if err != nil {
			return err
		}
		manager, err := app.Service(e, app.ServiceNotification)
		if err != nil {
			return err
		}
		return e.InvokeVoid(manager, "notify",
			jni.Sig(jni.TVoid, jni.TInt, jni.TClass("android/app/Notification")),
			jni.Int(n.ID), jni.Ref(built))
	})
}
