//go:build android

package notify

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// addAction builds one button and adds it to the notification.
func addAction(e *jni.Env, builder jni.Object, id, index int32, a Action, token int64) error {
	if a.Name == "" {
		return fmt.Errorf("notify: the action %q has no name", a.Title)
	}
	pending, err := broadcast(e, id, index, a.Name, token, a.Reply != nil)
	if err != nil {
		return err
	}
	title, err := e.String(a.Title)
	if err != nil {
		return err
	}
	ab, err := e.Make("android/app/Notification$Action$Builder",
		jni.Sig(jni.TVoid, jni.TInt, jni.TClass("java/lang/CharSequence"),
			jni.TClass("android/app/PendingIntent")),
		jni.Int(a.Icon), jni.Ref(title), jni.Ref(pending))
	if err != nil {
		return err
	}
	if a.Reply != nil {
		input, err := remoteInput(e, a.Name, a.Reply.Label)
		if err != nil {
			return err
		}
		if _, err := e.Invoke(ab, "addRemoteInput",
			jni.Sig(jni.TClass("android/app/Notification$Action$Builder"),
				jni.TClass("android/app/RemoteInput")), jni.Ref(input)); err != nil {
			return err
		}
	}
	built, err := e.Invoke(ab, "build",
		jni.Sig(jni.TClass("android/app/Notification$Action")))
	if err != nil {
		return err
	}
	_, err = e.Invoke(builder, "addAction",
		jni.Sig(jni.TClass("android/app/Notification$Builder"),
			jni.TClass("android/app/Notification$Action")), jni.Ref(built))
	return err
}
