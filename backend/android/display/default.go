//go:build android

package display

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Default is the screen the app is on.
func Default() (Display, error) {
	var out Display
	err := jni.Do(func(e *jni.Env) error {
		wm, err := app.Service(e, app.ServiceWindow)
		if err != nil {
			return err
		}
		d, err := e.Invoke(wm, "getDefaultDisplay",
			jni.Sig(jni.TClass("android/view/Display")))
		if err != nil {
			return err
		}
		if d.IsNil() {
			return errNoDisplay
		}
		out, err = read(e, d)
		return err
	})
	return out, err
}
