//go:build android

package notify

import "github.com/gabrielluizsf/antui/backend/android/jni"

// builderProgress sets a progress bar when there is one, and the silent flag
// when the platform understands it.
func builderProgress(e *jni.Env, builder jni.Object, n Notification, sdk int) error {
	b := jni.TClass("android/app/Notification$Builder")
	if n.Total > 0 || n.Indeterminate {
		if err := setBuilder(e, builder, "setProgress", jni.Sig(b, jni.TInt, jni.TInt, jni.TBool),
			jni.Int(n.Total), jni.Int(n.Progress), jni.Bool(n.Indeterminate)); err != nil {
			return err
		}
	}
	if n.Silent && sdk >= 26 {
		// Below 26 the sound is a property of the notification and above it
		// a property of the channel, so this only means anything one side of
		// that line. setSilent is API 29: not having it is not a failure.
		setBuilder(e, builder, "setSilent", jni.Sig(b, jni.TBool), jni.Bool(true))
	}
	return nil
}
