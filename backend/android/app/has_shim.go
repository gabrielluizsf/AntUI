//go:build android

package app

import "github.com/gabrielluizsf/antui/backend/android/jni"

// HasShim reports whether this app was built with the Java shim, and so
// whether permission results, activity results and arriving intents can be
// delivered at all.
func HasShim() bool {
	found := false
	jni.Do(func(e *jni.Env) error {
		_, err := e.Class(ShimClass)
		found = err == nil
		return nil
	})
	return found
}
