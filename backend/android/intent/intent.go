//go:build android

package intent

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// The actions this package uses, spelled as the platform spells them.
const (
	ActionView = "android.intent.action.VIEW"
	ActionSend = "android.intent.action.SEND"
	ActionEdit = "android.intent.action.EDIT"
)

// OpenURL hands a link to whatever the user opens links with — a browser, or
// the app that claims that address.
func OpenURL(url string) error {
	return jni.Do(func(e *jni.Env) error {
		in, err := viewIntent(e, url, "")
		if err != nil {
			return err
		}
		return start(e, in)
	})
}

// View opens something with the app that handles its type. mime may be empty,
// in which case the system guesses from the address.
func View(uri, mime string) error {
	return jni.Do(func(e *jni.Env) error {
		in, err := viewIntent(e, uri, mime)
		if err != nil {
			return err
		}
		return start(e, in)
	})
}

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

// grantRead is Intent.FLAG_GRANT_READ_URI_PERMISSION: the receiving app may
// read the address this intent carries, and only that address, and only
// until the app that sent it goes away.
const grantRead int32 = 0x00000001

// ShareFile offers a file to any app that takes its kind.
//
// The address has to be a content:// one from a provider — see
// [antui/backend/android/share], which is what makes one. A path will not do, and
// handing one over throws rather than failing quietly, which is the platform
// being firm about something that used to be a way to leak private files.
func ShareFile(uri, mime, title string) error {
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
		jmime, err := e.String(mime)
		if err != nil {
			return err
		}
		if _, err := e.Invoke(in, "setType",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString),
			jni.Ref(jmime)); err != nil {
			return err
		}
		juri, err := e.String(uri)
		if err != nil {
			return err
		}
		parsed, err := e.Static("android/net/Uri", "parse",
			jni.Sig(jni.TClass("android/net/Uri"), jni.TString), jni.Ref(juri))
		if err != nil {
			return err
		}
		// EXTRA_STREAM carries a Uri and not a string. Putting the text of
		// one in works as far as the compiler is concerned and produces a
		// share every app refuses.
		key, err := e.String("android.intent.extra.STREAM")
		if err != nil {
			return err
		}
		if _, err := e.Invoke(in, "putExtra",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString,
				jni.TClass("android/os/Parcelable")),
			jni.Ref(key), jni.Ref(parsed)); err != nil {
			return err
		}
		if _, err := e.Invoke(in, "addFlags",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TInt),
			jni.Int(grantRead)); err != nil {
			return err
		}
		chooser, err := wrapInChooser(e, in, title)
		if err != nil {
			return err
		}
		// The grant has to be on the chooser too: it is the chooser that is
		// started, and the flags on what it wraps do not travel.
		if _, err := e.Invoke(chooser, "addFlags",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TInt),
			jni.Int(grantRead)); err != nil {
			return err
		}
		return start(e, chooser)
	})
}

// Start sends an intent the caller built itself, for everything this package
// has no name for.
func Start(in jni.Object) error {
	return jni.Do(func(e *jni.Env) error { return start(e, in) })
}

// Reading what arrived.

// Action is what an intent asks for: [ActionView] for a link that was
// opened, and so on.
func Action(in jni.Object) (string, error) {
	if in.IsNil() {
		return "", nil
	}
	var out string
	err := jni.Do(func(e *jni.Env) error {
		var err error
		out, err = e.InvokeString(in, "getAction", jni.Sig(jni.TString))
		return err
	})
	return out, err
}

// Data is the address an intent carries — the link that was followed, the
// file that was opened. It is empty when there is none.
func Data(in jni.Object) (string, error) {
	if in.IsNil() {
		return "", nil
	}
	var out string
	err := jni.Do(func(e *jni.Env) error {
		var err error
		out, err = e.InvokeString(in, "getDataString", jni.Sig(jni.TString))
		return err
	})
	return out, err
}

// Extra reads one string an intent was given — which is how a notification
// says which notification it was.
func Extra(in jni.Object, key string) (string, error) {
	if in.IsNil() {
		return "", nil
	}
	var out string
	err := jni.Do(func(e *jni.Env) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		out, err = e.InvokeString(in, "getStringExtra",
			jni.Sig(jni.TString, jni.TString), jni.Ref(k))
		return err
	})
	return out, err
}

// Launch is the intent the app was started with, which is where a deep link
// arrives when the app was not already running. [app.OnNewIntent] is where
// it arrives when it was.
func Launch() (jni.Object, error) {
	var out jni.Object
	err := jni.Do(func(e *jni.Env) error {
		in, err := e.Invoke(app.Context(), "getIntent",
			jni.Sig(jni.TClass("android/content/Intent")))
		if err != nil {
			return err
		}
		if !in.IsNil() {
			out = e.Global(in)
		}
		return nil
	})
	return out, err
}

// The plumbing.

func viewIntent(e *jni.Env, address, mime string) (jni.Object, error) {
	action, err := e.String(ActionView)
	if err != nil {
		return jni.Object{}, err
	}
	jaddr, err := e.String(address)
	if err != nil {
		return jni.Object{}, err
	}
	uri, err := e.Static("android/net/Uri", "parse",
		jni.Sig(jni.TClass("android/net/Uri"), jni.TString), jni.Ref(jaddr))
	if err != nil {
		return jni.Object{}, err
	}
	in, err := e.Make("android/content/Intent",
		jni.Sig(jni.TVoid, jni.TString, jni.TClass("android/net/Uri")),
		jni.Ref(action), jni.Ref(uri))
	if err != nil {
		return jni.Object{}, err
	}
	if mime != "" {
		jmime, err := e.String(mime)
		if err != nil {
			return jni.Object{}, err
		}
		if _, err := e.Invoke(in, "setDataAndType",
			jni.Sig(jni.TClass("android/content/Intent"),
				jni.TClass("android/net/Uri"), jni.TString),
			jni.Ref(uri), jni.Ref(jmime)); err != nil {
			return jni.Object{}, err
		}
	}
	return in, nil
}

func putString(e *jni.Env, in jni.Object, key, value string) error {
	k, err := e.String(key)
	if err != nil {
		return err
	}
	v, err := e.String(value)
	if err != nil {
		return err
	}
	_, err = e.Invoke(in, "putExtra",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TString),
		jni.Ref(k), jni.Ref(v))
	return err
}

func wrapInChooser(e *jni.Env, in jni.Object, title string) (jni.Object, error) {
	jtitle, err := e.String(title)
	if err != nil {
		return jni.Object{}, err
	}
	return e.Static("android/content/Intent", "createChooser",
		jni.Sig(jni.TClass("android/content/Intent"),
			jni.TClass("android/content/Intent"), jni.TClass("java/lang/CharSequence")),
		jni.Ref(in), jni.Ref(jtitle))
}

// start hands the intent to the system, and turns the one failure everybody
// hits into a sentence rather than a Java class name.
func start(e *jni.Env, in jni.Object) error {
	err := e.InvokeVoid(app.Context(), "startActivity",
		jni.Sig(jni.TVoid, jni.TClass("android/content/Intent")), jni.Ref(in))
	var je *jni.Error
	if errorsAs(err, &je) && je.Class == "android.content.ActivityNotFoundException" {
		return fmt.Errorf("intent: nothing on this device can open that")
	}
	return err
}

func errorsAs(err error, target **jni.Error) bool {
	je, ok := err.(*jni.Error)
	if ok {
		*target = je
	}
	return ok
}
