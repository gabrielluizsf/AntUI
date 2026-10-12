//go:build android

package prefs

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// with runs f with the SharedPreferences object.
func (p *Prefs) with(f func(*jni.Env, jni.Object) error) error {
	return jni.Do(func(e *jni.Env) error {
		var obj jni.Object
		var err error
		if p.name == "" {
			obj, err = e.Static("android/preference/PreferenceManager",
				"getDefaultSharedPreferences",
				jni.Sig(jni.TClass("android/content/SharedPreferences"), jni.TContext),
				jni.Ref(app.Context()))
		} else {
			var jname jni.Object
			if jname, err = e.String(p.name); err != nil {
				return err
			}
			obj, err = e.Invoke(app.Context(), "getSharedPreferences",
				jni.Sig(jni.TClass("android/content/SharedPreferences"), jni.TString, jni.TInt),
				jni.Ref(jname), jni.Int(modePrivate))
		}
		if err != nil {
			return err
		}
		return f(e, obj)
	})
}
