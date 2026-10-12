//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// makeChannel builds the NotificationChannel and registers it.
func makeChannel(e *jni.Env, c Channel) error {
	id, err := e.String(c.ID)
	if err != nil {
		return err
	}
	name, err := e.String(c.Name)
	if err != nil {
		return err
	}
	ch, err := e.Make("android/app/NotificationChannel",
		jni.Sig(jni.TVoid, jni.TString, jni.TClass("java/lang/CharSequence"), jni.TInt),
		jni.Ref(id), jni.Ref(name), jni.Int(int32(c.Importance)))
	if err != nil {
		return err
	}
	if c.Description != "" {
		d, err := e.String(c.Description)
		if err != nil {
			return err
		}
		if err := e.InvokeVoid(ch, "setDescription",
			jni.Sig(jni.TVoid, jni.TString), jni.Ref(d)); err != nil {
			return err
		}
	}
	manager, err := app.Service(e, app.ServiceNotification)
	if err != nil {
		return err
	}
	return e.InvokeVoid(manager, "createNotificationChannel",
		jni.Sig(jni.TVoid, jni.TClass("android/app/NotificationChannel")), jni.Ref(ch))
}
