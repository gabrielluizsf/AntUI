//go:build android

package power

import "github.com/gabrielluizsf/antui/backend/android/jni"

// readExtras reads the values out of the battery intent.
func readExtras(e *jni.Env, status jni.Object, out *Battery) error {
	extra := func(name string) (int32, error) {
		k, err := e.String(name)
		if err != nil {
			return 0, err
		}
		return e.InvokeInt(status, "getIntExtra",
			jni.Sig(jni.TInt, jni.TString, jni.TInt), jni.Ref(k), jni.Int(-1))
	}
	level, err := extra("level")
	if err != nil {
		return err
	}
	scale, err := extra("scale")
	if err != nil {
		return err
	}
	if level >= 0 && scale > 0 {
		out.Level = float32(level) / float32(scale)
	}
	state, err := extra("status")
	if err != nil {
		return err
	}
	out.Charging = state == statusCharging
	out.Full = state == statusFull
	plugged, err := extra("plugged")
	if err != nil {
		return err
	}
	if plugged > 0 {
		out.Plugged = int(plugged)
	}
	return nil
}
