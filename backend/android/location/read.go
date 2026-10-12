//go:build android

package location

import (
	"time"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// read pulls a Fix out of an android.location.Location.
func read(e *jni.Env, loc jni.Object) (Fix, error) {
	var f Fix
	var err error
	if f.Latitude, err = invokeDouble(e, loc, "getLatitude"); err != nil {
		return f, err
	}
	if f.Longitude, err = invokeDouble(e, loc, "getLongitude"); err != nil {
		return f, err
	}
	if f.Altitude, err = invokeDouble(e, loc, "getAltitude"); err != nil {
		return f, err
	}
	if f.Accuracy, err = e.InvokeFloat(loc, "getAccuracy", jni.Sig(jni.TFloat)); err != nil {
		return f, err
	}
	if f.Speed, err = e.InvokeFloat(loc, "getSpeed", jni.Sig(jni.TFloat)); err != nil {
		return f, err
	}
	if f.Bearing, err = e.InvokeFloat(loc, "getBearing", jni.Sig(jni.TFloat)); err != nil {
		return f, err
	}
	ms, err := e.InvokeLong(loc, "getTime", jni.Sig(jni.TLong))
	if err != nil {
		return f, err
	}
	f.At = time.UnixMilli(ms)
	if f.Provider, err = e.InvokeString(loc, "getProvider", jni.Sig(jni.TString)); err != nil {
		return f, err
	}
	return f, nil
}
