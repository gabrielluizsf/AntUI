//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Allowed reports whether the user has left notifications on for this app.
// An app that posts while this is false is posting into nothing.
func Allowed() (bool, error) {
	var out bool
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServiceNotification)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(manager, "areNotificationsEnabled", jni.Sig(jni.TBool))
		return err
	})
	return out, err
}
