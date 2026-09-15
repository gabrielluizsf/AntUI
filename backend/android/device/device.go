//go:build android

// Package device is what the machine says about itself.
package device

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Info is the device, as far as an app is allowed to know it.
type Info struct {
	// Manufacturer, Brand and Model are what is printed on the box and in
	// the settings screen: "Google", "google", "Pixel 8".
	Manufacturer string
	Brand        string
	Model        string
	// Device and Product are the internal names, which is what a bug report
	// and a compatibility list are keyed on.
	Device  string
	Product string
	// SDK is the API level, and Release is the version a person would say:
	// 34 and "14".
	SDK     int
	Release string
	// Locale is the language and region the user chose, as "pt-BR".
	Locale string
	// TimeZone is the zone's id, as "America/Sao_Paulo".
	TimeZone string
	// Emulator is a guess, from the build's own fingerprint. It is right far
	// more often than not and should never be the basis of anything but a
	// diagnostic.
	Emulator bool
}

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

// ID is an identifier for this app on this device.
//
// It is Settings.Secure.ANDROID_ID, which since Android 8 is **different for
// every app** and is reset when the app is uninstalled or the device is
// wiped. That is deliberate on the platform's part: there is no longer any
// identifier that follows a person across apps, and an app that wants one
// anyway is one the store will refuse.
func ID() (string, error) {
	var id string
	err := jni.Do(func(e *jni.Env) error {
		resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver")))
		if err != nil {
			return err
		}
		name, err := e.String("android_id")
		if err != nil {
			return err
		}
		id, err = func() (string, error) {
			o, err := e.Static("android/provider/Settings$Secure", "getString",
				jni.Sig(jni.TString, jni.TClass("android/content/ContentResolver"), jni.TString),
				jni.Ref(resolver), jni.Ref(name))
			if err != nil {
				return "", err
			}
			return e.GoString(o)
		}()
		return err
	})
	return id, err
}

// looksEmulated is a guess and is named like one.
func looksEmulated(fingerprint, model, manufacturer string) bool {
	for _, sign := range []string{"generic", "unknown", "emulator", "sdk_gphone",
		"Android SDK built for", "Genymotion", "vbox"} {
		if contains(fingerprint, sign) || contains(model, sign) || contains(manufacturer, sign) {
			return true
		}
	}
	return false
}

func contains(haystack, needle string) bool {
	if len(needle) > len(haystack) {
		return false
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if equalFold(haystack[i:i+len(needle)], needle) {
			return true
		}
	}
	return false
}

func equalFold(a, b string) bool {
	for i := range len(a) {
		x, y := a[i], b[i]
		if x >= 'A' && x <= 'Z' {
			x += 32
		}
		if y >= 'A' && y <= 'Z' {
			y += 32
		}
		if x != y {
			return false
		}
	}
	return true
}
