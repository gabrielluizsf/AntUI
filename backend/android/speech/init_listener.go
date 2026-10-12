//go:build android

package speech

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// initListener reports the engine's onInit status once, and drops a second
// one rather than blocking — a channel of one is all that is needed.
func initListener(ready chan int32) app.Handler {
	return func(e *jni.Env, method string, args []jni.Object) {
		if method != "onInit" || len(args) == 0 {
			return
		}
		status, err := e.InvokeInt(args[0], "intValue", jni.Sig(jni.TInt))
		if err != nil {
			status = -1
		}
		select {
		case ready <- status:
		default:
		}
	}
}
