//go:build android

package clipboard

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Text is what is on the clipboard, or the empty string when there is
// nothing, or nothing that is text.
//
// Two things the platform does that this cannot hide. Since Android 10 an
// app may only read the clipboard **while it has focus** — a background read
// returns nothing rather than failing, which is the platform closing a way
// apps used to spy on each other. And since Android 12 a read puts a toast
// on the screen saying the app pasted, which is not something an app can
// turn off.
func Text() (string, error) {
	var out string
	var err error
	if e := app.RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			manager, err := app.Service(e, app.ServiceClipboard)
			if err != nil {
				return err
			}
			has, err := e.InvokeBool(manager, "hasPrimaryClip", jni.Sig(jni.TBool))
			if err != nil || !has {
				return err
			}
			clip, err := e.Invoke(manager, "getPrimaryClip",
				jni.Sig(jni.TClass("android/content/ClipData")))
			if err != nil || clip.IsNil() {
				return err
			}
			n, err := e.InvokeInt(clip, "getItemCount", jni.Sig(jni.TInt))
			if err != nil || n == 0 {
				return err
			}
			item, err := e.Invoke(clip, "getItemAt",
				jni.Sig(jni.TClass("android/content/ClipData$Item"), jni.TInt), jni.Int(0))
			if err != nil || item.IsNil() {
				return err
			}
			// getText gives a CharSequence, which may be styled text rather
			// than a String, so it goes through toString rather than being
			// read directly.
			text, err := e.Invoke(item, "getText",
				jni.Sig(jni.TClass("java/lang/CharSequence")))
			if err != nil || text.IsNil() {
				return err
			}
			out, err = e.InvokeString(text, "toString", jni.Sig(jni.TString))
			return err
		})
	}); e != nil {
		return "", e
	}
	if err == nil && out == "" {
		if a := app.Current(); a != nil && !a.Focused() {
			return "", ErrNotFocused
		}
	}
	return out, err
}
