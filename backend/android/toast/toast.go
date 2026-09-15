//go:build android

// Package toast is the small message that appears over an app and fades.
//
// It is the platform's own, not something drawn into the canvas — which
// means it survives the app being covered, looks like every other toast on
// the device, and cannot be styled. A message that belongs to the app's own
// picture belongs in the app's own picture.
package toast

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// How long a toast stays. The platform allows no other lengths.
const (
	Short = 0 // about two seconds
	Long  = 1 // about three and a half
)

// Show puts a message on the screen for a moment.
//
// From Android 12 a toast from an app that is not in front is refused by the
// system and nothing appears. There is no error: the call succeeds and the
// message is dropped.
func Show(text string, length int) error {
	var err error
	if e := app.RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			jtext, err := e.String(text)
			if err != nil {
				return err
			}
			t, err := e.Static("android/widget/Toast", "makeText",
				jni.Sig(jni.TClass("android/widget/Toast"),
					jni.TContext, jni.TClass("java/lang/CharSequence"), jni.TInt),
				jni.Ref(app.Context()), jni.Ref(jtext), jni.Int(int32(length)))
			if err != nil {
				return err
			}
			return e.InvokeVoid(t, "show", jni.Sig(jni.TVoid))
		})
	}); e != nil {
		return e
	}
	return err
}
