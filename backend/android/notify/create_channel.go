//go:build android

package notify

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// CreateChannel makes a channel, or does nothing if it already exists.
// Calling it every launch is the intended use: it is cheap and it is how a
// name or a description gets updated.
//
// Below API 26 there are no channels and this does nothing successfully.
func CreateChannel(c Channel) error {
	if c.ID == "" {
		return fmt.Errorf("notify: a channel needs an id")
	}
	a := app.Current()
	if a == nil {
		return app.ErrNoUIThread
	}
	if a.Config().SDK < 26 {
		return nil
	}
	if c.Importance == 0 {
		c.Importance = Default
	}
	return jni.Do(func(e *jni.Env) error {
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
	})
}
