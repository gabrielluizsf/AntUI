//go:build android

package picker

import "github.com/gabrielluizsf/antui/backend/android/jni"

// drop gives back the global reference the activity result came with. It is
// the caller's to free, and this package is that caller.
func drop(o jni.Object) {
	if o.IsNil() {
		return
	}
	jni.Do(func(e *jni.Env) error {
		e.DeleteGlobal(o)
		return nil
	})
}
