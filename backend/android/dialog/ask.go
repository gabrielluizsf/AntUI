//go:build android

package dialog

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

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
