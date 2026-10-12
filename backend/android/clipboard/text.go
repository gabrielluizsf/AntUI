//go:build android

package clipboard

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Text is what is on the clipboard, or the empty string when there is
// nothing, or nothing that is text.
//
// Two things the platform does that this cannot hide. Since Android 10 an
// app may only read the clipboard **while it has focus** — a background read
// returns nothing rather than failing, which is the platform closing a way
// apps used to spy on each other. And since Android 12 a read puts a toast
// on the screen saying the app pasted, which is not something an app can
// turn off.
func Text() (string, error) {
	var out string
	var err error
	if e := app.RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			out, err = readText(e)
			return err
		})
	}); e != nil {
		return "", e
	}
	if err == nil && out == "" {
		if a := app.Current(); a != nil && !a.Focused() {
			return "", ErrNotFocused
		}
	}
	return out, err
}
