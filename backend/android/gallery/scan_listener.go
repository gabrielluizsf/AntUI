//go:build android

package gallery

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// scanListener reports the address the scanner gave the file, or empty when
// it would not take it.
func scanListener(answer chan string) app.Handler {
	return func(e *jni.Env, method string, args []jni.Object) {
		if method != "onScanCompleted" {
			return
		}
		out := ""
		if len(args) > 1 && !args[1].IsNil() {
			out, _ = e.InvokeString(args[1], "toString", jni.Sig(jni.TString))
		}
		select {
		case answer <- out:
		default:
		}
	}
}
