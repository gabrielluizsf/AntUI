//go:build android

package display

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Rotation is how far the screen is turned from the way the device was
// built: 0, 1, 2 or 3 for none, a quarter, a half and three quarters
// clockwise. It is what [antui/backend/android/sensor.Reading.ForDisplay] takes,
// and the reason it has to be asked for at all is that the sensors do not
// turn with the screen.
func Rotation() (int, error) {
	var out int
	err := jni.Do(func(e *jni.Env) error {
		wm, err := app.Service(e, app.ServiceWindow)
		if err != nil {
			return err
		}
		d, err := e.Invoke(wm, "getDefaultDisplay",
			jni.Sig(jni.TClass("android/view/Display")))
		if err != nil || d.IsNil() {
			return err
		}
		r, err := e.InvokeInt(d, "getRotation", jni.Sig(jni.TInt))
		out = int(r)
		return err
	})
	return out, err
}
