//go:build android

package notify

import "github.com/gabrielluizsf/antui/backend/android/jni"

func setAction(e *jni.Env, intent jni.Object, action string) error {
	js, err := e.String(action)
	if err != nil {
		return err
	}
	_, err = e.Invoke(intent, "setAction",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString), jni.Ref(js))
	return err
}
