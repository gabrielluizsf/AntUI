//go:build android

package power

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Saving reports whether the user has turned battery saver on. An app that
// respects it does less: fewer frames, no background work, no animation
// that is only decoration.
func Saving() (bool, error) {
	var out bool
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServicePower)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(manager, "isPowerSaveMode", jni.Sig(jni.TBool))
		return err
	})
	return out, err
}
