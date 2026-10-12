//go:build android

package app

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Service is a system service by name — the object every one of Android's
// managers is got from.
func Service(e *jni.Env, name string) (jni.Object, error) {
	js, err := e.String(name)
	if err != nil {
		return jni.Object{}, err
	}
	svc, err := e.Invoke(Context(), "getSystemService",
		jni.Sig(jni.TObject, jni.TString), jni.Ref(js))
	if err != nil {
		return jni.Object{}, err
	}
	if svc.IsNil() {
		return jni.Object{}, fmt.Errorf("app: this device has no %q service", name)
	}
	return svc, nil
}
