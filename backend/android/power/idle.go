//go:build android

package power

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Idle reports whether the device is in doze — the deep sleep it enters when
// it has been still and unplugged for a while, in which network access and
// alarms are held until it wakes.
//
// It is false below API 23, where there was no doze.
func Idle() (bool, error) {
	a := app.Current()
	if a == nil {
		return false, app.ErrNoUIThread
	}
	if a.Config().SDK < 23 {
		return false, nil
	}
	var out bool
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServicePower)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(manager, "isDeviceIdleMode", jni.Sig(jni.TBool))
		return err
	})
	return out, err
}
