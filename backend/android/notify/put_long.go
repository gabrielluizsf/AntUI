//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func putLong(e *jni.Env, intent jni.Object, key string, value int64) error {
	k, err := e.String(key)
	if err != nil {
		return err
	}
	_, err = e.Invoke(intent, "putExtra",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TLong),
		jni.Ref(k), jni.Long(value))
	return err
}
