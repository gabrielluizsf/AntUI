//go:build android

package biometric

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func fingerprintAvailability() (Availability, error) {
	var out Availability
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, "fingerprint")
		if err != nil {
			return err
		}
		hardware, err := e.InvokeBool(manager, "isHardwareDetected", jni.Sig(jni.TBool))
		if err != nil {
			return err
		}
		if !hardware {
			out = NoHardware
			return nil
		}
		enrolled, err := e.InvokeBool(manager, "hasEnrolledFingerprints", jni.Sig(jni.TBool))
		if err != nil {
			return err
		}
		if !enrolled {
			out = NoneEnrolled
			return nil
		}
		out = Available
		return nil
	})
	return out, err
}
