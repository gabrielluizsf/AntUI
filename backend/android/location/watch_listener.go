//go:build android

package location

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// watchListener reports each fix, and quietly ignores anything malformed.
func watchListener(f func(Fix)) app.Handler {
	return func(e *jni.Env, method string, args []jni.Object) {
		if method != "onLocationChanged" || len(args) == 0 || args[0].IsNil() {
			return
		}
		if fix, err := read(e, args[0]); err == nil {
			f(fix)
		}
	}
}
