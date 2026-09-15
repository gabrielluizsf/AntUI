//go:build android

package picker

import (
	"errors"
	"os"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// The intent flags that decide what may be done with an address and for how
// long.
const (
	grantRead  int32 = 0x00000001
	grantWrite int32 = 0x00000002
	// grantPersist is what makes a permission survive a restart. Without it
	// the address the user chose works until the app is closed and then
	// stops, which looks like the app losing the file.
	grantPersist int32 = 0x00000040
)

// Open reads a file the user chose.
//
// The address is a content:// one, not a path, and it may have no file
// behind it at all — the document could be in a cloud account, being
// streamed as it is read. So this is the only way to get at one, and the
// file it gives back is a descriptor rather than a name on disk.
func Open(uri string) (*os.File, error) { return openMode(uri, "r") }

// Writer opens a chosen file for writing, replacing what is in it.
func Writer(uri string) (*os.File, error) { return openMode(uri, "wt") }

func openMode(uri, mode string) (*os.File, error) {
	var file *os.File
	err := jni.Do(func(e *jni.Env) error {
		resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver")))
		if err != nil {
			return err
		}
		parsed, err := parse(e, uri)
		if err != nil {
			return err
		}
		jmode, err := e.String(mode)
		if err != nil {
			return err
		}
		pfd, err := e.Invoke(resolver, "openFileDescriptor",
			jni.Sig(jni.TClass("android/os/ParcelFileDescriptor"),
				jni.TClass("android/net/Uri"), jni.TString),
			jni.Ref(parsed), jni.Ref(jmode))
		if err != nil {
			return err
		}
		if pfd.IsNil() {
			return errors.New("picker: nothing behind that address")
		}
		fd, err := e.InvokeInt(pfd, "detachFd", jni.Sig(jni.TInt))
		if err != nil {
			return err
		}
		file = os.NewFile(uintptr(fd), uri)
		if file == nil {
			return errors.New("picker: the descriptor is not usable")
		}
		return nil
	})
	return file, err
}

// Name is what the user would call the file — the name shown in the picker,
// which is not in the address and has to be asked for.
func Name(uri string) (string, error) {
	var out string
	err := jni.Do(func(e *jni.Env) error {
		resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver")))
		if err != nil {
			return err
		}
		parsed, err := parse(e, uri)
		if err != nil {
			return err
		}
		cursor, err := e.Invoke(resolver, "query",
			jni.Sig(jni.TClass("android/database/Cursor"),
				jni.TClass("android/net/Uri"), jni.TArray(jni.TString),
				jni.TString, jni.TArray(jni.TString), jni.TString),
			jni.Ref(parsed), jni.Ref(jni.Object{}), jni.Ref(jni.Object{}),
			jni.Ref(jni.Object{}), jni.Ref(jni.Object{}))
		if err != nil || cursor.IsNil() {
			return err
		}
		defer e.InvokeVoid(cursor, "close", jni.Sig(jni.TVoid))

		moved, err := e.InvokeBool(cursor, "moveToFirst", jni.Sig(jni.TBool))
		if err != nil || !moved {
			return err
		}
		col, err := e.String("_display_name")
		if err != nil {
			return err
		}
		at, err := e.InvokeInt(cursor, "getColumnIndex",
			jni.Sig(jni.TInt, jni.TString), jni.Ref(col))
		if err != nil || at < 0 {
			return err
		}
		out, err = e.InvokeString(cursor, "getString",
			jni.Sig(jni.TString, jni.TInt), jni.Int(at))
		return err
	})
	return out, err
}

// Create asks the user where to put a new file, and gives back its address.
//
// name is what the picker suggests; the user may change it, and what comes
// back is whatever they settled on. Nothing is written — the file exists and
// is empty, and [Writer] is how it is filled.
func Create(name, mime string) (string, error) {
	var in jni.Object
	err := jni.Do(func(e *jni.Env) error {
		var err error
		in, err = document(e, "android.intent.action.CREATE_DOCUMENT", mime, name, false)
		return err
	})
	if err != nil {
		return "", err
	}
	return single(in)
}

// Folder asks the user for a whole directory, and gives back its address.
//
// It is how an app is given a place to work in rather than one file — a
// folder of photographs to process, a directory to save into every time. The
// address is a *tree*, and it is not something [Open] takes: reading what is
// in it means the document API, which this library does not wrap yet.
//
// Pair it with [Keep], or the choice is forgotten the moment the app closes.
func Folder() (string, error) {
	var in jni.Object
	err := jni.Do(func(e *jni.Env) error {
		action, err := e.String("android.intent.action.OPEN_DOCUMENT_TREE")
		if err != nil {
			return err
		}
		made, err := e.Make("android/content/Intent",
			jni.Sig(jni.TVoid, jni.TString), jni.Ref(action))
		if err != nil {
			return err
		}
		if _, err := e.Invoke(made, "addFlags",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TInt),
			jni.Int(grantRead|grantWrite|grantPersist)); err != nil {
			return err
		}
		in = e.Global(made)
		return nil
	})
	if err != nil {
		return "", err
	}
	return single(in)
}

// single starts an intent that answers with one address.
func single(in jni.Object) (string, error) {
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
		uri, err := e.Invoke(r.Data, "getData", jni.Sig(jni.TClass("android/net/Uri")))
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

// Keep asks the system to remember that the user chose this address, so that
// it still works after the app has been closed and opened again.
//
// Without it a chosen file works until the process ends and then stops,
// which looks to a user exactly like the app losing their document. There is
// a limit on how many an app may keep — a few hundred — and the oldest are
// dropped, so keeping every file a user ever opened is not a plan.
func Keep(uri string, write bool) error {
	flags := grantRead
	if write {
		flags |= grantWrite
	}
	return jni.Do(func(e *jni.Env) error {
		resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver")))
		if err != nil {
			return err
		}
		parsed, err := parse(e, uri)
		if err != nil {
			return err
		}
		return e.InvokeVoid(resolver, "takePersistableUriPermission",
			jni.Sig(jni.TVoid, jni.TClass("android/net/Uri"), jni.TInt),
			jni.Ref(parsed), jni.Int(flags))
	})
}

// Forget gives a kept permission back.
func Forget(uri string, write bool) error {
	flags := grantRead
	if write {
		flags |= grantWrite
	}
	return jni.Do(func(e *jni.Env) error {
		resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver")))
		if err != nil {
			return err
		}
		parsed, err := parse(e, uri)
		if err != nil {
			return err
		}
		return e.InvokeVoid(resolver, "releasePersistableUriPermission",
			jni.Sig(jni.TVoid, jni.TClass("android/net/Uri"), jni.TInt),
			jni.Ref(parsed), jni.Int(flags))
	})
}

// Kept is every address the app is still allowed to reach, from choices the
// user made in earlier runs.
func Kept() ([]string, error) {
	var out []string
	err := jni.Do(func(e *jni.Env) error {
		resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver")))
		if err != nil {
			return err
		}
		list, err := e.Invoke(resolver, "getPersistedUriPermissions",
			jni.Sig(jni.TClass("java/util/List")))
		if err != nil || list.IsNil() {
			return err
		}
		n, err := e.InvokeInt(list, "size", jni.Sig(jni.TInt))
		if err != nil || n == 0 {
			return err
		}
		return e.Frame(int(n)*4+16, func() error {
			for i := range int(n) {
				item, err := e.Invoke(list, "get", jni.Sig(jni.TObject, jni.TInt),
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
	})
	return out, err
}

func parse(e *jni.Env, uri string) (jni.Object, error) {
	js, err := e.String(uri)
	if err != nil {
		return jni.Object{}, err
	}
	return e.Static("android/net/Uri", "parse",
		jni.Sig(jni.TClass("android/net/Uri"), jni.TString), jni.Ref(js))
}
