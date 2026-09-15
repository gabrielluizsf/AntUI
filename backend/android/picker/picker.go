//go:build android

package picker

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// ErrCancelled is the user backing out without choosing.
var ErrCancelled = errors.New("picker: the user chose nothing")

// Image asks the user to choose a picture, and gives back its address.
//
// **No permission is needed, and none should be asked for.** The user
// choosing a file *is* the permission: what comes back is a content://
// address this app may read and nothing else. Asking for READ_EXTERNAL_STORAGE
// to open a picker is asking for the whole library in order to be handed one
// photo, and the store now refuses apps that do it.
//
// The address is not a path. It is opened with
// [antui/backend/android/app.Context]'s content resolver; there is no file behind it
// on many devices, because the picture may live in a cloud account.
//
// **Not from the frame loop.** The picker is another app, in front, for as
// long as the user wants it there.
func Image() (string, error) { return one("image/*") }

// Video is Image for a film.
func Video() (string, error) { return one("video/*") }

// Audio is Image for a sound.
func Audio() (string, error) { return one("audio/*") }

// Any asks for a file of a given type — "application/pdf", "text/*", or
// "*/*" for anything at all.
func Any(mime string) (string, error) { return one(mime) }

// Many asks for several at once. It comes back with one address per file, in
// the order the picker gave them.
func Many(mime string) ([]string, error) {
	in, err := build(mime, true)
	if err != nil {
		return nil, err
	}
	r, err := app.StartForResult(in)
	if err != nil {
		return nil, err
	}
	defer drop(r.Data)
	if !r.OK() {
		return nil, ErrCancelled
	}
	return addresses(r.Data)
}

func one(mime string) (string, error) {
	in, err := build(mime, false)
	if err != nil {
		return "", err
	}
	r, err := app.StartForResult(in)
	if err != nil {
		return "", err
	}
	defer drop(r.Data)
	if !r.OK() {
		return "", ErrCancelled
	}
	var out string
	err = jni.Do(func(e *jni.Env) error {
		uri, err := e.Invoke(r.Data, "getData",
			jni.Sig(jni.TClass("android/net/Uri")))
		if err != nil || uri.IsNil() {
			return err
		}
		out, err = e.InvokeString(uri, "toString", jni.Sig(jni.TString))
		return err
	})
	if err != nil {
		return "", err
	}
	if out == "" {
		return "", ErrCancelled
	}
	return out, nil
}

// document builds an OPEN_DOCUMENT or CREATE_DOCUMENT intent.
func document(e *jni.Env, action, mime, title string, many bool) (jni.Object, error) {
	jaction, err := e.String(action)
	if err != nil {
		return jni.Object{}, err
	}
	in, err := e.Make("android/content/Intent",
		jni.Sig(jni.TVoid, jni.TString), jni.Ref(jaction))
	if err != nil {
		return jni.Object{}, err
	}
	jmime, err := e.String(mime)
	if err != nil {
		return jni.Object{}, err
	}
	if _, err := e.Invoke(in, "setType",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString), jni.Ref(jmime)); err != nil {
		return jni.Object{}, err
	}
	cat, err := e.String("android.intent.category.OPENABLE")
	if err != nil {
		return jni.Object{}, err
	}
	if _, err := e.Invoke(in, "addCategory",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString), jni.Ref(cat)); err != nil {
		return jni.Object{}, err
	}
	if title != "" {
		key, err := e.String("android.intent.extra.TITLE")
		if err != nil {
			return jni.Object{}, err
		}
		jtitle, err := e.String(title)
		if err != nil {
			return jni.Object{}, err
		}
		if _, err := e.Invoke(in, "putExtra",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TString),
			jni.Ref(key), jni.Ref(jtitle)); err != nil {
			return jni.Object{}, err
		}
	}
	if _, err := e.Invoke(in, "addFlags",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TInt),
		jni.Int(grantRead|grantWrite|grantPersist)); err != nil {
		return jni.Object{}, err
	}
	return e.Global(in), nil
}

// build makes the intent.
//
// ACTION_OPEN_DOCUMENT rather than ACTION_GET_CONTENT: the first goes to the
// system's own picker and gives an address that keeps working, and the second
// goes to whichever app claims the type and gives one that may not outlive
// the call. The system picker is also the one that shows cloud accounts.
func build(mime string, many bool) (jni.Object, error) {
	var out jni.Object
	err := jni.Do(func(e *jni.Env) error {
		action, err := e.String("android.intent.action.OPEN_DOCUMENT")
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
		// CATEGORY_OPENABLE keeps out the things that have no bytes behind
		// them, which the picker will otherwise happily offer.
		cat, err := e.String("android.intent.category.OPENABLE")
		if err != nil {
			return err
		}
		if _, err := e.Invoke(in, "addCategory",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString),
			jni.Ref(cat)); err != nil {
			return err
		}
		// Read, and persistable so that Keep can make the choice survive a
		// restart. A picker started without the persist flag gives an
		// address that stops working when the process does, and takes the
		// permission with it.
		if _, err := e.Invoke(in, "addFlags",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TInt),
			jni.Int(grantRead|grantPersist)); err != nil {
			return err
		}
		if many {
			key, err := e.String("android.intent.extra.ALLOW_MULTIPLE")
			if err != nil {
				return err
			}
			if _, err := e.Invoke(in, "putExtra",
				jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TBool),
				jni.Ref(key), jni.Bool(true)); err != nil {
				return err
			}
		}
		out = e.Global(in)
		return nil
	})
	return out, err
}

// addresses reads the one address or the list of them, because the picker
// returns them in two different places depending on how many were chosen —
// even when multiple were allowed and the user picked one.
func addresses(data jni.Object) ([]string, error) {
	var out []string
	err := jni.Do(func(e *jni.Env) error {
		clip, err := e.Invoke(data, "getClipData",
			jni.Sig(jni.TClass("android/content/ClipData")))
		if err != nil {
			return err
		}
		if !clip.IsNil() {
			n, err := e.InvokeInt(clip, "getItemCount", jni.Sig(jni.TInt))
			if err != nil {
				return err
			}
			return e.Frame(int(n)*4+16, func() error {
				for i := range int(n) {
					item, err := e.Invoke(clip, "getItemAt",
						jni.Sig(jni.TClass("android/content/ClipData$Item"), jni.TInt),
						jni.Int(int32(i)))
					if err != nil || item.IsNil() {
						continue
					}
					uri, err := e.Invoke(item, "getUri",
						jni.Sig(jni.TClass("android/net/Uri")))
					if err != nil || uri.IsNil() {
						continue
					}
					s, err := e.InvokeString(uri, "toString", jni.Sig(jni.TString))
					if err == nil && s != "" {
						out = append(out, s)
					}
				}
				return nil
			})
		}
		uri, err := e.Invoke(data, "getData", jni.Sig(jni.TClass("android/net/Uri")))
		if err != nil || uri.IsNil() {
			return err
		}
		s, err := e.InvokeString(uri, "toString", jni.Sig(jni.TString))
		if err == nil && s != "" {
			out = append(out, s)
		}
		return err
	})
	return out, err
}

// drop gives back the global reference the activity result came with. It is
// the caller's to free, and this package is that caller.
func drop(o jni.Object) {
	if o.IsNil() {
		return
	}
	jni.Do(func(e *jni.Env) error {
		e.DeleteGlobal(o)
		return nil
	})
}
