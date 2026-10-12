//go:build android

package clipboard

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// HasText reports whether there is anything on the clipboard. It is subject
// to the same focus rule as [Text].
func HasText() (bool, error) {
	var out bool
	var err error
	if e := app.RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			manager, err := app.Service(e, app.ServiceClipboard)
			if err != nil {
				return err
			}
			out, err = e.InvokeBool(manager, "hasPrimaryClip", jni.Sig(jni.TBool))
			return err
		})
	}); e != nil {
		return false, e
	}
	return out, err
}
