//go:build android

package permission

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// ShouldExplain reports whether the user has refused this permission once
// and the app should say why it wants it before asking again.
//
// It is a three-state answer squeezed into two, and the third state is the
// one that matters: **false means either "never asked" or "refused for
// good"**. The platform does not distinguish them, deliberately — an app
// that could tell the difference could pester. So the way to use it is: ask;
// if refused and this is true, explain and ask once more; if refused and
// this is false, the dialog will not appear again and the only way forward
// is the settings screen.
func ShouldExplain(name string) (bool, error) {
	a := app.Current()
	if a == nil {
		return false, app.ErrNoUIThread
	}
	if a.Config().SDK < 23 {
		return false, nil
	}
	var out bool
	err := jni.Do(func(e *jni.Env) error {
		js, err := e.String(name)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(app.Context(), "shouldShowRequestPermissionRationale",
			jni.Sig(jni.TBool, jni.TString), jni.Ref(js))
		return err
	})
	return out, err
}
