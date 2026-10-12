//go:build android

package audio

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// unpluggedHandler calls f on the becoming-noisy broadcast.
func unpluggedHandler(f func()) app.Handler {
	return func(e *jni.Env, method string, args []jni.Object) {
		if method == "onReceive" && f != nil {
			f()
		}
	}
}
