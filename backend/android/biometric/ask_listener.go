//go:build android

package biometric

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// askListener turns the prompt's callbacks into a single answer.
func askListener(finish func(error)) app.Handler {
	return func(e *jni.Env, method string, args []jni.Object) {
		switch method {
		case "onAuthenticationSucceeded":
			finish(nil)
		case "onAuthenticationFailed":
			// Not an answer: a finger the sensor did not recognise. The
			// prompt is still up.
		case "onAuthenticationError":
			finish(promptError(e, args))
		case "onClick":
			// The negative button, which arrives here as well as through an
			// error on most versions.
			finish(ErrRefused)
		}
	}
}
