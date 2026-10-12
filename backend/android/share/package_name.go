//go:build android

package share

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func packageName() (string, error) {
	var err error
	pkgOnce.Do(func() {
		err = jni.Do(func(e *jni.Env) error {
			var e2 error
			pkgName, e2 = e.InvokeString(app.Context(), "getPackageName", jni.Sig(jni.TString))
			return e2
		})
	})
	if err != nil {
		return "", err
	}
	if pkgName == "" {
		return "", errors.New("share: the app does not know its own name")
	}
	return pkgName, nil
}
