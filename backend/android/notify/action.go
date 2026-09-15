//go:build android

package notify

import (
	"fmt"
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Action is a button on a notification.
//
// It carries a **name**, not a function, and that is not an inconvenience —
// it is the only thing that works. A notification outlives the app that
// posted it: the user may tap the button an hour later, on a device that
// killed the process long ago, and Android starts the app afresh to deliver
// it. A Go closure cannot survive that. A name can, and [OnAction] is where
// the app says what each name means, at startup, every time it starts.
type Action struct {
	// Title is the label on the button.
	Title string
	// Name is what [OnAction] matched it against. It never appears on screen.
	Name string
	// Icon is a drawable resource id, and may be zero: from Android 7 the
	// icon on an action button is not drawn at all, and it is only still
	// required because the API predates that.
	Icon int32
	// Reply turns the button into a text field the user can type into, which
	// is how a message is answered from the shade.
	Reply *Reply
}

// Reply is a text field on an action.
type Reply struct {
	// Label is the hint shown in the empty field.
	Label string
}

// The intent extras an action travels in.
const (
	extraAction = "dev.antui.action"
	extraID     = "dev.antui.id"
)

var (
	actionMu   sync.Mutex
	actionFns  = map[string]func(id int32, reply string){}
	actionTok  int64
	actionOnce sync.Once
)

// OnAction says what a named action does.
//
// Register every one of them at startup, before anything is posted, and
// register them again on every launch — the tap that brings the app back may
// be the first thing that happens in a new process.
//
// f runs on a thread of the system's choosing, not the app's. reply is what
// the user typed, and is empty for a button that is only a button.
func OnAction(name string, f func(id int32, reply string)) {
	actionMu.Lock()
	actionFns[name] = f
	actionMu.Unlock()
	actionOnce.Do(startReceiver)
}

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

// addActions puts the buttons on a notification being built.
func addActions(e *jni.Env, builder jni.Object, n Notification) error {
	if len(n.Actions) == 0 {
		return nil
	}
	actionMu.Lock()
	token := actionTok
	actionMu.Unlock()
	if token == 0 {
		return fmt.Errorf("notify: this notification has actions and none has " +
			"been registered with OnAction, so a tap would go nowhere")
	}

	for i, a := range n.Actions {
		if a.Name == "" {
			return fmt.Errorf("notify: the action %q has no name", a.Title)
		}
		pending, err := broadcast(e, n.ID, int32(i), a.Name, token, a.Reply != nil)
		if err != nil {
			return err
		}
		title, err := e.String(a.Title)
		if err != nil {
			return err
		}
		ab, err := e.Make("android/app/Notification$Action$Builder",
			jni.Sig(jni.TVoid, jni.TInt, jni.TClass("java/lang/CharSequence"),
				jni.TClass("android/app/PendingIntent")),
			jni.Int(a.Icon), jni.Ref(title), jni.Ref(pending))
		if err != nil {
			return err
		}
		if a.Reply != nil {
			input, err := remoteInput(e, a.Name, a.Reply.Label)
			if err != nil {
				return err
			}
			if _, err := e.Invoke(ab, "addRemoteInput",
				jni.Sig(jni.TClass("android/app/Notification$Action$Builder"),
					jni.TClass("android/app/RemoteInput")), jni.Ref(input)); err != nil {
				return err
			}
		}
		built, err := e.Invoke(ab, "build",
			jni.Sig(jni.TClass("android/app/Notification$Action")))
		if err != nil {
			return err
		}
		if _, err := e.Invoke(builder, "addAction",
			jni.Sig(jni.TClass("android/app/Notification$Builder"),
				jni.TClass("android/app/Notification$Action")), jni.Ref(built)); err != nil {
			return err
		}
	}
	return nil
}

// broadcast builds the PendingIntent one action sends.
//
// The mutable flag is the trap. A PendingIntent has had to be immutable since
// API 31 — except one carrying a reply, because the system's whole job there
// is to *add* the typed text to it. Marking a reply immutable produces a
// button that works and a reply that is always empty.
func broadcast(e *jni.Env, id, index int32, name string, token int64, mutable bool) (jni.Object, error) {
	cls, err := e.Class("dev/antui/AntuiReceiver")
	if err != nil {
		return jni.Object{}, err
	}
	intent, err := e.Make("android/content/Intent",
		jni.Sig(jni.TVoid, jni.TContext, jni.TClass("java/lang/Class")),
		jni.Ref(app.Context()), jni.Ref(cls.Object()))
	if err != nil {
		return jni.Object{}, err
	}
	// A distinct action string per button, so that two PendingIntents for
	// the same notification are not treated as the same one and merged —
	// which they are if only their extras differ.
	if err := setAction(e, intent, fmt.Sprintf("dev.antui.%d.%d", id, index)); err != nil {
		return jni.Object{}, err
	}
	if err := putString(e, intent, extraAction, name); err != nil {
		return jni.Object{}, err
	}
	if err := putInt(e, intent, extraID, id); err != nil {
		return jni.Object{}, err
	}
	if err := putLong(e, intent, app.ReceiverToken, token); err != nil {
		return jni.Object{}, err
	}

	flags := flagUpdateCurrent | flagImmutable
	if mutable {
		flags = flagUpdateCurrent | flagMutable
	}
	// The request code has to differ per button for the same reason the
	// action string does.
	return e.Static("android/app/PendingIntent", "getBroadcast",
		jni.Sig(jni.TClass("android/app/PendingIntent"),
			jni.TContext, jni.TInt, jni.TClass("android/content/Intent"), jni.TInt),
		jni.Ref(app.Context()), jni.Int(id*100+index), jni.Ref(intent), jni.Int(flags))
}

func remoteInput(e *jni.Env, key, label string) (jni.Object, error) {
	jkey, err := e.String(key)
	if err != nil {
		return jni.Object{}, err
	}
	b, err := e.Make("android/app/RemoteInput$Builder",
		jni.Sig(jni.TVoid, jni.TString), jni.Ref(jkey))
	if err != nil {
		return jni.Object{}, err
	}
	if label != "" {
		jlabel, err := e.String(label)
		if err != nil {
			return jni.Object{}, err
		}
		if _, err := e.Invoke(b, "setLabel",
			jni.Sig(jni.TClass("android/app/RemoteInput$Builder"),
				jni.TClass("java/lang/CharSequence")), jni.Ref(jlabel)); err != nil {
			return jni.Object{}, err
		}
	}
	return e.Invoke(b, "build", jni.Sig(jni.TClass("android/app/RemoteInput")))
}

// flagMutable lets the system add to the intent, which a reply needs.
const flagMutable int32 = 1 << 25

func setAction(e *jni.Env, intent jni.Object, action string) error {
	js, err := e.String(action)
	if err != nil {
		return err
	}
	_, err = e.Invoke(intent, "setAction",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString), jni.Ref(js))
	return err
}

func putString(e *jni.Env, intent jni.Object, key, value string) error {
	k, err := e.String(key)
	if err != nil {
		return err
	}
	v, err := e.String(value)
	if err != nil {
		return err
	}
	_, err = e.Invoke(intent, "putExtra",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TString),
		jni.Ref(k), jni.Ref(v))
	return err
}

func putInt(e *jni.Env, intent jni.Object, key string, value int32) error {
	k, err := e.String(key)
	if err != nil {
		return err
	}
	_, err = e.Invoke(intent, "putExtra",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TInt),
		jni.Ref(k), jni.Int(value))
	return err
}

func putLong(e *jni.Env, intent jni.Object, key string, value int64) error {
	k, err := e.String(key)
	if err != nil {
		return err
	}
	_, err = e.Invoke(intent, "putExtra",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TLong),
		jni.Ref(k), jni.Long(value))
	return err
}

func stringExtra(e *jni.Env, intent jni.Object, key string) (string, error) {
	k, err := e.String(key)
	if err != nil {
		return "", err
	}
	return e.InvokeString(intent, "getStringExtra",
		jni.Sig(jni.TString, jni.TString), jni.Ref(k))
}

func mustString(e *jni.Env, s string) jni.Value {
	o, err := e.String(s)
	if err != nil {
		return jni.Ref(jni.Object{})
	}
	return jni.Ref(o)
}
