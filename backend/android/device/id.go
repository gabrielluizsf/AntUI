//go:build android

package device

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

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
