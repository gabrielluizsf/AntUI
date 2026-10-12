//go:build android

package location

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Enabled reports whether a provider is switched on. The permission being
// granted is not enough: the user can turn location off for the whole
// device, and then every request answers nothing for ever.
func Enabled(provider string) (bool, error) {
	var out bool
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServiceLocation)
		if err != nil {
			return err
		}
		jp, err := e.String(provider)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(manager, "isProviderEnabled",
			jni.Sig(jni.TBool, jni.TString), jni.Ref(jp))
		return err
	})
	return out, err
}
