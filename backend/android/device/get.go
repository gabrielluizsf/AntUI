//go:build android

package device

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Get reads everything at once, which is one connection to the machine
// rather than a dozen.
func Get() (Info, error) {
	var out Info
	err := jni.Do(func(e *jni.Env) error {
		for _, f := range []struct {
			name string
			into *string
		}{
			{"MANUFACTURER", &out.Manufacturer},
			{"BRAND", &out.Brand},
			{"MODEL", &out.Model},
			{"DEVICE", &out.Device},
			{"PRODUCT", &out.Product},
		} {
			v, err := e.ConstantString("android/os/Build", f.name)
			if err != nil {
				return err
			}
			*f.into = v
		}
		sdk, err := e.ConstantInt("android/os/Build$VERSION", "SDK_INT")
		if err != nil {
			return err
		}
		out.SDK = int(sdk)
		if out.Release, err = e.ConstantString("android/os/Build$VERSION", "RELEASE"); err != nil {
			return err
		}

		zone, err := e.Static("java/util/TimeZone", "getDefault",
			jni.Sig(jni.TClass("java/util/TimeZone")))
		if err != nil {
			return err
		}
		if out.TimeZone, err = e.InvokeString(zone, "getID", jni.Sig(jni.TString)); err != nil {
			return err
		}

		fingerprint, err := e.ConstantString("android/os/Build", "FINGERPRINT")
		if err != nil {
			return err
		}
		out.Emulator = looksEmulated(fingerprint, out.Model, out.Manufacturer)
		return nil
	})
	if a := app.Current(); a != nil {
		out.Locale = a.Config().Locale()
	}
	return out, err
}
