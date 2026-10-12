//go:build android

package picker

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// single starts an intent that answers with one address.
func single(in jni.Object) (string, error) {
	r, err := app.StartForResult(in)
	if err != nil {
		return "", err
	}
	defer drop(r.Data)
	if !r.OK() {
		return "", ErrCancelled
	}
	var out string
	err = jni.Do(func(e *jni.Env) error {
		uri, err := e.Invoke(r.Data, "getData", jni.Sig(jni.TClass("android/net/Uri")))
		if err != nil || uri.IsNil() {
			return err
		}
		out, err = e.InvokeString(uri, "toString", jni.Sig(jni.TString))
		return err
	})
	if err != nil {
		return "", err
	}
	if out == "" {
		return "", ErrCancelled
	}
	return out, nil
}
