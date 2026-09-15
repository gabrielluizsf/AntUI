//go:build android

package dialog

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Which button was pressed. The numbers are the platform's, and they are
// negative — DialogInterface counts its buttons downwards from -1.
type Button int

const (
	// Cancelled is the back button or a tap outside, which is not one of the
	// buttons and is the answer most code forgets to handle.
	Cancelled Button = 0
	// Positive is the one on the right, which is the one that does the
	// thing: OK, Save, Delete.
	Positive Button = -1
	// Negative is the one that does not: Cancel.
	Negative Button = -2
	// Neutral is the third, for an answer that is neither — "Not now",
	// "Tell me more".
	Neutral Button = -3
)

// String names the button, which is what a log line wants.
func (b Button) String() string {
	switch b {
	case Positive:
		return "positive"
	case Negative:
		return "negative"
	case Neutral:
		return "neutral"
	}
	return "cancelled"
}

// Ask puts the platform's own dialog on screen and waits for an answer.
//
// It is the system's dialog, not one drawn into the canvas: it looks like
// every other dialog on the device, it survives the app being covered, and
// it cannot be styled. **A dialog drawn by antui itself is the better answer
// for most apps** — it matches the rest of what is on screen and needs no
// round trip through Java. This is for the times when looking like the system
// is the point: a confirmation before deleting, a permission being explained.
//
// A button with an empty label is left out. Cancelling — the back button, or
// a tap outside — answers [Cancelled].
//
// **Not from the frame loop.** It waits for a person, and while it waits the
// app is paused and Android is waiting for the app to say so.
func Ask(title, message string, positive, negative, neutral string) (Button, error) {
	answer := make(chan Button, 1)
	var once sync.Once
	send := func(b Button) { once.Do(func() { answer <- b }) }

	token, release := app.Listen(func(e *jni.Env, method string, args []jni.Object) {
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
	})
	defer release()

	var err error
	if e := app.RunOnUISync(func() {
		err = show(token, title, message, positive, negative, neutral)
	}); e != nil {
		return Cancelled, e
	}
	if err != nil {
		return Cancelled, err
	}
	return <-answer, nil
}

// Message is a dialog with one button, for saying something rather than
// asking. It waits until the button is pressed.
func Message(title, message, dismiss string) error {
	if dismiss == "" {
		dismiss = "OK"
	}
	_, err := Ask(title, message, dismiss, "", "")
	return err
}

// Confirm is a yes-or-no dialog, and reports whether the answer was yes.
func Confirm(title, message, yes, no string) (bool, error) {
	b, err := Ask(title, message, yes, no, "")
	return b == Positive, err
}

// The interfaces the dialog answers through. All three are interfaces, so
// one proxy implements the lot and no Java had to be written for any of them.
var listenerInterfaces = []string{
	"android.content.DialogInterface$OnClickListener",
	"android.content.DialogInterface$OnCancelListener",
	"android.content.DialogInterface$OnDismissListener",
}

var builderType = jni.TClass("android/app/AlertDialog$Builder")

func show(token int64, title, message, positive, negative, neutral string) error {
	return jni.Do(func(e *jni.Env) error {
		builder, err := e.Make("android/app/AlertDialog$Builder",
			jni.Sig(jni.TVoid, jni.TContext), jni.Ref(app.Context()))
		if err != nil {
			return err
		}
		listener, err := app.Proxy(e, token, listenerInterfaces...)
		if err != nil {
			return err
		}

		text := func(method, value string) error {
			if value == "" {
				return nil
			}
			js, err := e.String(value)
			if err != nil {
				return err
			}
			_, err = e.Invoke(builder, method,
				jni.Sig(builderType, jni.TClass("java/lang/CharSequence")), jni.Ref(js))
			return err
		}
		if err := text("setTitle", title); err != nil {
			return err
		}
		if err := text("setMessage", message); err != nil {
			return err
		}
		for _, b := range []struct{ method, label string }{
			{"setPositiveButton", positive},
			{"setNegativeButton", negative},
			{"setNeutralButton", neutral},
		} {
			if b.label == "" {
				continue
			}
			js, err := e.String(b.label)
			if err != nil {
				return err
			}
			if _, err := e.Invoke(builder, b.method,
				jni.Sig(builderType, jni.TClass("java/lang/CharSequence"),
					jni.TClass("android/content/DialogInterface$OnClickListener")),
				jni.Ref(js), jni.Ref(listener)); err != nil {
				return err
			}
		}
		// The back button and a tap outside both have to answer, or a caller
		// waiting on the channel waits forever.
		if _, err := e.Invoke(builder, "setOnCancelListener",
			jni.Sig(builderType,
				jni.TClass("android/content/DialogInterface$OnCancelListener")),
			jni.Ref(listener)); err != nil {
			return err
		}
		_, err = e.Invoke(builder, "show",
			jni.Sig(jni.TClass("android/app/AlertDialog")))
		return err
	})
}
