//go:build android

package biometric

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Can reports whether a prompt would work.
func Can() (Availability, error) {
	a := app.Current()
	if a == nil {
		return TooOld, app.ErrNoUIThread
	}
	sdk := a.Config().SDK
	if sdk < 28 {
		return TooOld, nil
	}
	if sdk < 29 {
		// BiometricManager arrived one version after the prompt did, so on
		// exactly this version the only way to ask is the older fingerprint
		// API — which is what "there is a sensor and something enrolled on
		// it" meant before faces existed.
		return fingerprintAvailability()
	}
	var out Availability
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, "biometric")
		if err != nil {
			return err
		}
		code, err := e.InvokeInt(manager, "canAuthenticate", jni.Sig(jni.TInt))
		if err != nil {
			return err
		}
		switch code {
		case 0:
			out = Available
		case 11:
			out = NoneEnrolled
		case 12:
			out = NoHardware
		default:
			out = Unavailable
		}
		return nil
	})
	return out, plainly(err)
}
