//go:build android

package intent

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Start sends an intent the caller built itself, for everything this package
// has no name for.
func Start(in jni.Object) error {
	return jni.Do(func(e *jni.Env) error { return start(e, in) })
}
