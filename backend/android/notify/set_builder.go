//go:build android

package notify

import "github.com/gabrielluizsf/antui/backend/android/jni"

// setBuilder calls one setter on the builder and drops its return, which is
// what a builder's setters give back.
func setBuilder(e *jni.Env, builder jni.Object, method, sig string, args ...jni.Value) error {
	_, err := e.Invoke(builder, method, sig, args...)
	return err
}
