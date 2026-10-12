//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// CancelAll takes all of this app's notifications back.
func CancelAll() error {
	return jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServiceNotification)
		if err != nil {
			return err
		}
		return e.InvokeVoid(manager, "cancelAll", jni.Sig(jni.TVoid))
	})
}
