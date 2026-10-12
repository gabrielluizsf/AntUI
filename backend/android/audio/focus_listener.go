//go:build android

package audio

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// focusListener reports each change of focus.
func focusListener(onChange func(Change)) app.Handler {
	return func(e *jni.Env, method string, args []jni.Object) {
		if method != "onAudioFocusChange" || len(args) == 0 {
			return
		}
		which, err := e.InvokeInt(args[0], "intValue", jni.Sig(jni.TInt))
		if err != nil {
			return
		}
		if onChange != nil {
			onChange(Change(which))
		}
	}
}
