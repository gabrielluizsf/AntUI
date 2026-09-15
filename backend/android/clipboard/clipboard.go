//go:build android

package clipboard

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// ErrNotFocused is a read refused because the app is not the one the user is
// looking at.
//
// The platform does not say so: it hands back an empty clipboard, exactly as
// if there were nothing on it, and writes the refusal to the system log
// where the app cannot see it. That silence is a day of somebody's life, so
// this package looks at whether it has focus and says which of the two
// happened.
var ErrNotFocused = errors.New("clipboard: the app is not in focus, and " +
	"Android only lets the app in focus read the clipboard")

// SetText puts text on the clipboard.
//
// label is what the system shows when it describes what was copied — some
// versions put it in a toast. It is not the text and is not pasted anywhere.
func SetText(label, text string) error {
	var err error
	if e := app.RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			manager, err := app.Service(e, app.ServiceClipboard)
			if err != nil {
				return err
			}
			jlabel, err := e.String(label)
			if err != nil {
				return err
			}
			jtext, err := e.String(text)
			if err != nil {
				return err
			}
			// A ClipData is a label plus one or more items, because the
			// clipboard can hold a URI or an intent as easily as a string.
			clip, err := e.Static("android/content/ClipData", "newPlainText",
				jni.Sig(jni.TClass("android/content/ClipData"),
					jni.TClass("java/lang/CharSequence"),
					jni.TClass("java/lang/CharSequence")),
				jni.Ref(jlabel), jni.Ref(jtext))
			if err != nil {
				return err
			}
			return e.InvokeVoid(manager, "setPrimaryClip",
				jni.Sig(jni.TVoid, jni.TClass("android/content/ClipData")), jni.Ref(clip))
		})
	}); e != nil {
		return e
	}
	return err
}

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

// HasText reports whether there is anything on the clipboard. It is subject
// to the same focus rule as [Text].
func HasText() (bool, error) {
	var out bool
	var err error
	if e := app.RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			manager, err := app.Service(e, app.ServiceClipboard)
			if err != nil {
				return err
			}
			out, err = e.InvokeBool(manager, "hasPrimaryClip", jni.Sig(jni.TBool))
			return err
		})
	}); e != nil {
		return false, e
	}
	return out, err
}
