//go:build android

package gallery

import "github.com/gabrielluizsf/antui/backend/android/jni"

func putString(e *jni.Env, values jni.Object, key, value string) error {
	k, err := e.String(key)
	if err != nil {
		return err
	}
	v, err := e.String(value)
	if err != nil {
		return err
	}
	return e.InvokeVoid(values, "put",
		jni.Sig(jni.TVoid, jni.TString, jni.TString), jni.Ref(k), jni.Ref(v))
}
