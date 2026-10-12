//go:build android

package notify

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// addActions puts the buttons on a notification being built.
func addActions(e *jni.Env, builder jni.Object, n Notification) error {
	if len(n.Actions) == 0 {
		return nil
	}
	actionMu.Lock()
	token := actionTok
	actionMu.Unlock()
	if token == 0 {
		return fmt.Errorf("notify: this notification has actions and none has " +
			"been registered with OnAction, so a tap would go nowhere")
	}

	for i, a := range n.Actions {
		if a.Name == "" {
			return fmt.Errorf("notify: the action %q has no name", a.Title)
		}
		pending, err := broadcast(e, n.ID, int32(i), a.Name, token, a.Reply != nil)
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
		if _, err := e.Invoke(builder, "addAction",
			jni.Sig(jni.TClass("android/app/Notification$Builder"),
				jni.TClass("android/app/Notification$Action")), jni.Ref(built)); err != nil {
			return err
		}
	}
	return nil
}
