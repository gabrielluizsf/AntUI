//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// replyText reads what the user typed, which the system puts into the intent
// on its way to the receiver rather than into the notification.
func replyText(e *jni.Env, intent jni.Object, key string) string {
	bundle, err := e.Static("android/app/RemoteInput", "getResultsFromIntent",
		jni.Sig(jni.TClass("android/os/Bundle"), jni.TClass("android/content/Intent")),
		jni.Ref(intent))
	if err != nil || bundle.IsNil() {
		return ""
	}
	jkey, err := e.String(key)
	if err != nil {
		return ""
	}
	text, err := e.Invoke(bundle, "getCharSequence",
		jni.Sig(jni.TClass("java/lang/CharSequence"), jni.TString), jni.Ref(jkey))
	if err != nil || text.IsNil() {
		return ""
	}
	out, err := e.InvokeString(text, "toString", jni.Sig(jni.TString))
	if err != nil {
		return ""
	}
	return out
}
