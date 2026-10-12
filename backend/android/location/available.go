//go:build android

package location

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Available reports whether geocoding will work at all.
func Available() bool {
	var out bool
	jni.Do(func(e *jni.Env) error {
		var err error
		c, err := e.Class("android/location/Geocoder")
		if err != nil {
			return err
		}
		m, err := e.StaticMethod(c, "isPresent", jni.Sig(jni.TBool))
		if err != nil {
			return err
		}
		out, err = e.CallStaticBool(c, m)
		return err
	})
	return out
}
