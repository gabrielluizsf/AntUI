//go:build android

package intent

import "github.com/gabrielluizsf/antui/backend/android/jni"

func putString(e *jni.Env, in jni.Object, key, value string) error {
	k, err := e.String(key)
	if err != nil {
		return err
	}
	v, err := e.String(value)
	if err != nil {
		return err
	}
	_, err = e.Invoke(in, "putExtra",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TString),
		jni.Ref(k), jni.Ref(v))
	return err
}
