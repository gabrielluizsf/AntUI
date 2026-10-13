//go:build android

package notify

import "github.com/gabrielluizsf/antui/backend/android/jni"

// builderIntent sets what a tap opens, when there is anything.
func builderIntent(e *jni.Env, builder jni.Object, n Notification) error {
	pending, err := tapIntent(e, n.ID, n.Extra)
	if err != nil {
		return err
	}
	if pending.IsNil() {
		return nil
	}
	return setBuilder(e, builder, "setContentIntent",
		jni.Sig(jni.TClass("android/app/Notification$Builder"),
			jni.TClass("android/app/PendingIntent")), jni.Ref(pending))
}
