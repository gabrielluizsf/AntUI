//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// startReceiver registers the one handler every action button comes back
// through. There is one for the whole package, because the broadcast is
// matched to its action by a name in the intent rather than by which
// handler it was sent to.
func startReceiver() {
	token, _ := app.Listen(func(e *jni.Env, method string, args []jni.Object) {
		if method != "onReceive" || len(args) == 0 || args[0].IsNil() {
			return
		}
		intent := args[0]
		name, err := stringExtra(e, intent, extraAction)
		if err != nil || name == "" {
			return
		}
		id, err := e.InvokeInt(intent, "getIntExtra",
			jni.Sig(jni.TInt, jni.TString, jni.TInt), mustString(e, extraID), jni.Int(0))
		if err != nil {
			return
		}
		reply := replyText(e, intent, name)

		actionMu.Lock()
		f := actionFns[name]
		actionMu.Unlock()
		if f != nil {
			f(id, reply)
		}
	})
	actionMu.Lock()
	actionTok = token
	actionMu.Unlock()
}
