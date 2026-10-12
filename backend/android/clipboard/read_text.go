//go:build android

package clipboard

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// readText reads the first item of the primary clip.
func readText(e *jni.Env) (string, error) {
	manager, err := app.Service(e, app.ServiceClipboard)
	if err != nil {
		return "", err
	}
	has, err := e.InvokeBool(manager, "hasPrimaryClip", jni.Sig(jni.TBool))
	if err != nil || !has {
		return "", err
	}
	clip, err := e.Invoke(manager, "getPrimaryClip",
		jni.Sig(jni.TClass("android/content/ClipData")))
	if err != nil || clip.IsNil() {
		return "", err
	}
	n, err := e.InvokeInt(clip, "getItemCount", jni.Sig(jni.TInt))
	if err != nil || n == 0 {
		return "", err
	}
	item, err := e.Invoke(clip, "getItemAt",
		jni.Sig(jni.TClass("android/content/ClipData$Item"), jni.TInt), jni.Int(0))
	if err != nil || item.IsNil() {
		return "", err
	}
	// getText gives a CharSequence, which may be styled text rather
	// than a String, so it goes through toString rather than being
	// read directly.
	text, err := e.Invoke(item, "getText",
		jni.Sig(jni.TClass("java/lang/CharSequence")))
	if err != nil || text.IsNil() {
		return "", err
	}
	return e.InvokeString(text, "toString", jni.Sig(jni.TString))
}
