//go:build android

package notify

import "github.com/gabrielluizsf/antui/backend/android/jni"

func putInt(e *jni.Env, intent jni.Object, key string, value int32) error {
	k, err := e.String(key)
	if err != nil {
		return err
	}
	_, err = e.Invoke(intent, "putExtra",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TInt),
		jni.Ref(k), jni.Int(value))
	return err
}
