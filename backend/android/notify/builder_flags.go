//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// builderFlags sets the two boolean properties a notification has.
func builderFlags(e *jni.Env, builder jni.Object, n Notification) error {
	b := jni.TClass("android/app/Notification$Builder")
	if n.AutoCancel {
		if err := setBuilder(e, builder, "setAutoCancel", jni.Sig(b, jni.TBool), jni.Bool(true)); err != nil {
			return err
		}
	}
	if n.Ongoing {
		return setBuilder(e, builder, "setOngoing", jni.Sig(b, jni.TBool), jni.Bool(true))
	}
	return nil
}
