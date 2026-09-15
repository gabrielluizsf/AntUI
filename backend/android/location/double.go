//go:build android

package location

import "github.com/gabrielluizsf/antui/backend/android/jni"

// invokeDouble is the one shape jni.Invoke does not have a wrapper for,
// because a double return is rare everywhere except here — where three of the
// six numbers in a position are one.
func invokeDouble(e *jni.Env, o jni.Object, name string) (float64, error) {
	cls, err := e.ClassOf(o)
	if err != nil {
		return 0, err
	}
	m, err := e.Method(cls, name, jni.Sig(jni.TDouble))
	if err != nil {
		return 0, err
	}
	return e.CallDouble(o, m)
}
