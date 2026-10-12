//go:build android

package permission

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// OpenSettings takes the user to this app's settings page, which is the only
// way back from a permission refused for good.
func OpenSettings() error {
	return jni.Do(func(e *jni.Env) error {
		action, err := e.ConstantString("android/provider/Settings",
			"ACTION_APPLICATION_DETAILS_SETTINGS")
		if err != nil {
			return err
		}
		pkg, err := e.InvokeString(app.Context(), "getPackageName", jni.Sig(jni.TString))
		if err != nil {
			return err
		}
		jaction, err := e.String(action)
		if err != nil {
			return err
		}
		juri, err := e.String("package:" + pkg)
		if err != nil {
			return err
		}
		uri, err := e.Static("android/net/Uri", "parse",
			jni.Sig(jni.TClass("android/net/Uri"), jni.TString), jni.Ref(juri))
		if err != nil {
			return err
		}
		intent, err := e.Make("android/content/Intent",
			jni.Sig(jni.TVoid, jni.TString, jni.TClass("android/net/Uri")),
			jni.Ref(jaction), jni.Ref(uri))
		if err != nil {
			return err
		}
		return e.InvokeVoid(app.Context(), "startActivity",
			jni.Sig(jni.TVoid, jni.TClass("android/content/Intent")), jni.Ref(intent))
	})
}
