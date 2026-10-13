//go:build android

package picker

import "github.com/gabrielluizsf/antui/backend/android/jni"

// allowMultiple asks the picker to accept more than one file.
func allowMultiple(e *jni.Env, in jni.Object) error {
	key, err := e.String("android.intent.extra.ALLOW_MULTIPLE")
	if err != nil {
		return err
	}
	_, err = e.Invoke(in, "putExtra",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TBool),
		jni.Ref(key), jni.Bool(true))
	return err
}
