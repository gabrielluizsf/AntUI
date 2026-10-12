//go:build android

package speech

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// onEvent is the recogniser's callbacks. It runs on the UI thread, which is
// where startListening has to be called from.
func (l *listenState) onEvent(e *jni.Env, method string, args []jni.Object) {
	switch method {
	case "onResults", "onPartialResults":
		if len(args) == 0 {
			return
		}
		heard, ok := readResults(e, args[0], method == "onPartialResults")
		if ok {
			l.report(heard, nil)
		}
		l.mu.Lock()
		again := l.opt.Continuous && !l.stopped && method == "onResults"
		l.mu.Unlock()
		if again {
			startListening(e, l.recogniser, l.intent)
		}
	case "onError":
		if len(args) == 0 {
			return
		}
		code, err := e.InvokeInt(args[0], "intValue", jni.Sig(jni.TInt))
		if err != nil {
			return
		}
		l.report(Heard{}, recogniserError(code))
	}
}
