//go:build android

package intent

import "github.com/gabrielluizsf/antui/backend/android/jni"

// SendText offers text to any app that takes text: a message, an email, a
// note, the clipboard of another device.
//
// title is what the chooser sheet is headed with. Passing an empty one still
// shows a chooser — always going through one is deliberate, because sending
// straight to whichever app the user picked once is how an app ends up
// posting to the wrong place.
func SendText(text, subject, title string) error {
	return jni.Do(func(e *jni.Env) error {
		action, err := e.String(ActionSend)
		if err != nil {
			return err
		}
		in, err := e.Make("android/content/Intent",
			jni.Sig(jni.TVoid, jni.TString), jni.Ref(action))
		if err != nil {
			return err
		}
		mime, err := e.String("text/plain")
		if err != nil {
			return err
		}
		if _, err := e.Invoke(in, "setType",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString), jni.Ref(mime)); err != nil {
			return err
		}
		if err := putString(e, in, "android.intent.extra.TEXT", text); err != nil {
			return err
		}
		if subject != "" {
			if err := putString(e, in, "android.intent.extra.SUBJECT", subject); err != nil {
				return err
			}
		}
		chooser, err := wrapInChooser(e, in, title)
		if err != nil {
			return err
		}
		return start(e, chooser)
	})
}
