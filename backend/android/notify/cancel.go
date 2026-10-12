//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Cancel takes one notification back.
func Cancel(id int32) error {
	return jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServiceNotification)
		if err != nil {
			return err
		}
		return e.InvokeVoid(manager, "cancel", jni.Sig(jni.TVoid, jni.TInt), jni.Int(id))
	})
}
