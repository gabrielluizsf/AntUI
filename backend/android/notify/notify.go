//go:build android

package notify

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Importance is how loudly a channel announces itself. The user can change
// it afterwards and the app cannot change it back — a channel's importance
// is a request made once, when the channel is created, and never again.
type Importance int32

// The importances, as NotificationManager numbers them.
const (
	// Min shows nothing in the status bar.
	Min Importance = 1
	// Low appears silently.
	Low Importance = 2
	// Default makes a sound.
	Default Importance = 3
	// High makes a sound and slides in over whatever is on screen. It is for
	// things the user has to see now, and using it for anything else is how
	// an app gets its notifications turned off.
	High Importance = 4
)

// Channel is a category of notification the user can silence on its own.
//
// Every notification from API 26 onwards belongs to one, and posting to a
// channel that was never created shows nothing at all — with no error, which
// is the single most common way notifications fail to appear.
type Channel struct {
	// ID is what notifications refer to. It never changes and never appears
	// on screen.
	ID string
	// Name is what the user sees in the settings screen, so it should read
	// like a category: "Messages", "Downloads".
	Name string
	// Description is the line under the name, and may be empty.
	Description string
	Importance  Importance
}

// CreateChannel makes a channel, or does nothing if it already exists.
// Calling it every launch is the intended use: it is cheap and it is how a
// name or a description gets updated.
//
// Below API 26 there are no channels and this does nothing successfully.
func CreateChannel(c Channel) error {
	if c.ID == "" {
		return fmt.Errorf("notify: a channel needs an id")
	}
	a := app.Current()
	if a == nil {
		return app.ErrNoUIThread
	}
	if a.Config().SDK < 26 {
		return nil
	}
	if c.Importance == 0 {
		c.Importance = Default
	}
	return jni.Do(func(e *jni.Env) error {
		id, err := e.String(c.ID)
		if err != nil {
			return err
		}
		name, err := e.String(c.Name)
		if err != nil {
			return err
		}
		ch, err := e.Make("android/app/NotificationChannel",
			jni.Sig(jni.TVoid, jni.TString, jni.TClass("java/lang/CharSequence"), jni.TInt),
			jni.Ref(id), jni.Ref(name), jni.Int(int32(c.Importance)))
		if err != nil {
			return err
		}
		if c.Description != "" {
			d, err := e.String(c.Description)
			if err != nil {
				return err
			}
			if err := e.InvokeVoid(ch, "setDescription",
				jni.Sig(jni.TVoid, jni.TString), jni.Ref(d)); err != nil {
				return err
			}
		}
		manager, err := app.Service(e, app.ServiceNotification)
		if err != nil {
			return err
		}
		return e.InvokeVoid(manager, "createNotificationChannel",
			jni.Sig(jni.TVoid, jni.TClass("android/app/NotificationChannel")), jni.Ref(ch))
	})
}

// Notification is one message in the shade.
type Notification struct {
	// ID identifies it within this app. Posting twice with the same one
	// replaces rather than adds, which is how a progress bar moves.
	ID int32
	// Channel is the id of the channel it belongs to.
	Channel string
	Title   string
	Text    string

	// Icon is a drawable resource id. Zero means the platform's own
	// information icon, which is there on every device — an app has no
	// resources of its own until it is built with some.
	Icon int32

	// Progress draws a bar. Total is the top of the range; a Total of 0 means
	// no bar, and Indeterminate makes it a moving stripe with no position.
	Progress      int32
	Total         int32
	Indeterminate bool

	// Ongoing marks it as something in progress: it cannot be swiped away.
	Ongoing bool
	// AutoCancel dismisses it when the user taps it, which is what almost
	// every notification should do.
	AutoCancel bool
	// Silent posts without a sound even on a channel that makes one.
	Silent bool

	// Actions are the buttons under the notification. See [Action]: each
	// carries a name rather than a function, and [OnAction] is where the
	// names are given meaning.
	Actions []Action

	// Extra is handed back through [antui/backend/android/app.OnNewIntent] when the
	// user taps it, which is how an app knows *which* notification was
	// tapped.
	Extra map[string]string
}

// Post puts a notification in the shade, or replaces the one with the same
// id.
//
// From API 33 this needs the POST_NOTIFICATIONS permission, and without it
// nothing appears and nothing fails. Ask for it with
// [antui/backend/android/permission.Ensure] before the first post.
func Post(n Notification) error {
	a := app.Current()
	if a == nil {
		return app.ErrNoUIThread
	}
	sdk := a.Config().SDK
	return jni.Do(func(e *jni.Env) error {
		builder, err := newBuilder(e, sdk, n.Channel)
		if err != nil {
			return err
		}
		set := func(method, sig string, args ...jni.Value) error {
			_, err := e.Invoke(builder, method, sig, args...)
			return err
		}
		title, err := e.String(n.Title)
		if err != nil {
			return err
		}
		text, err := e.String(n.Text)
		if err != nil {
			return err
		}
		cs := jni.TClass("java/lang/CharSequence")
		b := jni.TClass("android/app/Notification$Builder")
		if err := set("setContentTitle", jni.Sig(b, cs), jni.Ref(title)); err != nil {
			return err
		}
		if err := set("setContentText", jni.Sig(b, cs), jni.Ref(text)); err != nil {
			return err
		}

		icon := n.Icon
		if icon == 0 {
			// Any notification without a small icon is refused outright, so
			// there has to be a default, and the only one every device has
			// is the platform's own.
			if icon, err = e.ConstantInt("android/R$drawable", "ic_dialog_info"); err != nil {
				return err
			}
		}
		if err := set("setSmallIcon", jni.Sig(b, jni.TInt), jni.Int(icon)); err != nil {
			return err
		}
		if n.AutoCancel {
			if err := set("setAutoCancel", jni.Sig(b, jni.TBool), jni.Bool(true)); err != nil {
				return err
			}
		}
		if n.Ongoing {
			if err := set("setOngoing", jni.Sig(b, jni.TBool), jni.Bool(true)); err != nil {
				return err
			}
		}
		if n.Total > 0 || n.Indeterminate {
			if err := set("setProgress", jni.Sig(b, jni.TInt, jni.TInt, jni.TBool),
				jni.Int(n.Total), jni.Int(n.Progress), jni.Bool(n.Indeterminate)); err != nil {
				return err
			}
		}
		if n.Silent && sdk >= 26 {
			// Below 26 the sound is a property of the notification and above
			// it a property of the channel, so this only means anything one
			// side of that line.
			if err := set("setSilent", jni.Sig(b, jni.TBool), jni.Bool(true)); err != nil {
				// setSilent is API 29. Not having it is not a failure.
				_ = err
			}
		}

		if err := addActions(e, builder, n); err != nil {
			return err
		}

		pending, err := tapIntent(e, n.ID, n.Extra)
		if err != nil {
			return err
		}
		if !pending.IsNil() {
			if err := set("setContentIntent",
				jni.Sig(b, jni.TClass("android/app/PendingIntent")), jni.Ref(pending)); err != nil {
				return err
			}
		}

		built, err := e.Invoke(builder, "build",
			jni.Sig(jni.TClass("android/app/Notification")))
		if err != nil {
			return err
		}
		manager, err := app.Service(e, app.ServiceNotification)
		if err != nil {
			return err
		}
		return e.InvokeVoid(manager, "notify",
			jni.Sig(jni.TVoid, jni.TInt, jni.TClass("android/app/Notification")),
			jni.Int(n.ID), jni.Ref(built))
	})
}

// newBuilder makes the right Notification.Builder for this platform: the one
// that takes a channel from API 26, and the one that does not below it.
func newBuilder(e *jni.Env, sdk int, channel string) (jni.Object, error) {
	if sdk >= 26 {
		if channel == "" {
			return jni.Object{}, fmt.Errorf(
				"notify: this device needs a channel, and none was named")
		}
		id, err := e.String(channel)
		if err != nil {
			return jni.Object{}, err
		}
		return e.Make("android/app/Notification$Builder",
			jni.Sig(jni.TVoid, jni.TContext, jni.TString),
			jni.Ref(app.Context()), jni.Ref(id))
	}
	return e.Make("android/app/Notification$Builder",
		jni.Sig(jni.TVoid, jni.TContext), jni.Ref(app.Context()))
}

// The PendingIntent flags. IMMUTABLE has been required since API 31 and is
// harmless before it; without it the call throws on a modern device.
const (
	flagUpdateCurrent int32 = 1 << 27
	flagImmutable     int32 = 1 << 26
	// Bring the activity to the front rather than starting a second copy.
	flagSingleTop int32 = 0x20000000
	flagClearTop  int32 = 0x04000000
)

// tapIntent builds what happens when the notification is tapped: the app is
// brought to the front, and the intent reaches it through onNewIntent.
func tapIntent(e *jni.Env, id int32, extra map[string]string) (jni.Object, error) {
	cls, err := e.Class(app.ShimClass)
	if err != nil {
		// Built without the shim: there is nothing that can receive the tap,
		// so the notification is posted without one rather than not at all.
		return jni.Object{}, nil
	}
	intent, err := e.Make("android/content/Intent",
		jni.Sig(jni.TVoid, jni.TContext, jni.TClass("java/lang/Class")),
		jni.Ref(app.Context()), jni.Ref(cls.Object()))
	if err != nil {
		return jni.Object{}, err
	}
	if _, err := e.Invoke(intent, "setFlags",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TInt),
		jni.Int(flagSingleTop|flagClearTop)); err != nil {
		return jni.Object{}, err
	}
	for k, v := range extra {
		jk, err := e.String(k)
		if err != nil {
			return jni.Object{}, err
		}
		jv, err := e.String(v)
		if err != nil {
			return jni.Object{}, err
		}
		if _, err := e.Invoke(intent, "putExtra",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TString),
			jni.Ref(jk), jni.Ref(jv)); err != nil {
			return jni.Object{}, err
		}
	}
	return e.Static("android/app/PendingIntent", "getActivity",
		jni.Sig(jni.TClass("android/app/PendingIntent"),
			jni.TContext, jni.TInt, jni.TClass("android/content/Intent"), jni.TInt),
		jni.Ref(app.Context()), jni.Int(id), jni.Ref(intent),
		jni.Int(flagUpdateCurrent|flagImmutable))
}

// Cancel takes one notification back.
func Cancel(id int32) error {
	return jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServiceNotification)
		if err != nil {
			return err
		}
		return e.InvokeVoid(manager, "cancel", jni.Sig(jni.TVoid, jni.TInt), jni.Int(id))
	})
}

// CancelAll takes all of this app's notifications back.
func CancelAll() error {
	return jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServiceNotification)
		if err != nil {
			return err
		}
		return e.InvokeVoid(manager, "cancelAll", jni.Sig(jni.TVoid))
	})
}

// Allowed reports whether the user has left notifications on for this app.
// An app that posts while this is false is posting into nothing.
func Allowed() (bool, error) {
	var out bool
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServiceNotification)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(manager, "areNotificationsEnabled", jni.Sig(jni.TBool))
		return err
	})
	return out, err
}
