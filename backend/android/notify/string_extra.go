//go:build android

package notify

import "github.com/gabrielluizsf/antui/backend/android/jni"

func stringExtra(e *jni.Env, intent jni.Object, key string) (string, error) {
	k, err := e.String(key)
	if err != nil {
		return "", err
	}
	return e.InvokeString(intent, "getStringExtra",
		jni.Sig(jni.TString, jni.TString), jni.Ref(k))
}
