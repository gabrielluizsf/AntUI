//go:build android

package picker

import "github.com/gabrielluizsf/antui/backend/android/jni"

func parse(e *jni.Env, uri string) (jni.Object, error) {
	js, err := e.String(uri)
	if err != nil {
		return jni.Object{}, err
	}
	return e.Static("android/net/Uri", "parse",
		jni.Sig(jni.TClass("android/net/Uri"), jni.TString), jni.Ref(js))
}
