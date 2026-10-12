//go:build android

package dialog

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// askListener turns the dialog's callbacks into a single answer.
func askListener(send func(Button)) app.Handler {
	return func(e *jni.Env, method string, args []jni.Object) {
		switch method {
		case "onClick":
			// onClick(DialogInterface, int). The button is the second
			// argument, boxed, because a proxy's arguments are all objects.
			if len(args) < 2 {
				send(Cancelled)
				return
			}
			which, err := e.InvokeInt(args[1], "intValue", jni.Sig(jni.TInt))
			if err != nil {
				send(Cancelled)
				return
			}
			send(Button(which))
		case "onCancel", "onDismiss":
			send(Cancelled)
		}
	}
}
