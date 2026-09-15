//go:build android

package gallery

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/permission"
)

// Kind is which of the device's libraries something belongs in.
type Kind int

const (
	// Image is the user's photos, which is where a screenshot or a picture
	// the app made belongs.
	Image Kind = iota
	// Video is their films.
	Video
	// Audio is their music and recordings.
	Audio
)

// collection is the address of the library, and folder the directory inside
// shared storage it lives in.
func (k Kind) collection() (class, field string) {
	switch k {
	case Video:
		return "android/provider/MediaStore$Video$Media", "EXTERNAL_CONTENT_URI"
	case Audio:
		return "android/provider/MediaStore$Audio$Media", "EXTERNAL_CONTENT_URI"
	}
	return "android/provider/MediaStore$Images$Media", "EXTERNAL_CONTENT_URI"
}

func (k Kind) folder() string {
	switch k {
	case Video:
		return "Movies"
	case Audio:
		return "Music"
	}
	return "Pictures"
}

// String names the library, which is what an error message wants.
func (k Kind) String() string {
	switch k {
	case Video:
		return "video"
	case Audio:
		return "audio"
	}
	return "image"
}

// The MediaStore columns used here, by the names the platform gives them.
const (
	colName    = "_display_name"
	colMime    = "mime_type"
	colPath    = "relative_path" // API 29 and up
	colPending = "is_pending"    // API 29 and up
)

// ErrNoPermission is saving on a device old enough to need one.
var ErrNoPermission = errors.New("gallery: below Android 10 an app needs " +
	"WRITE_EXTERNAL_STORAGE to put anything in the device's library")

// Writer is something being written into the device's library.
//
// Nothing is visible to the gallery until [Writer.Close]: on Android 10 and
// up the entry is marked pending while it is being written, so a half-copied
// video never appears in anyone's photo roll. Closing is what publishes it,
// and a Writer that is dropped without closing leaves an entry nothing can
// see and nothing will clean up — so close it, on the failure path too.
type Writer struct {
	file *os.File
	uri  jni.Object
	// uriText is the content:// address, which is what other apps are given.
	uriText string

	once   sync.Once
	closed bool
	err    error
}

// Create makes an entry in the device's library and gives back somewhere to
// write it.
//
// album is a folder inside the device's Pictures, Movies or Music — usually
// the app's name, so that what it saves is together and recognisable. It may
// be empty, and then the file goes in the top of the library.
//
// **On Android 10 and up this needs no permission at all.** An app owns what
// it creates: it may write it, read it back and delete it, and needs nothing
// granted to do so. Below that it needs WRITE_EXTERNAL_STORAGE, which is why
// this reports [ErrNoPermission] rather than failing at the write.
func Create(kind Kind, name, mime, album string) (*Writer, error) {
	a := app.Current()
	if a == nil {
		return nil, app.ErrNoUIThread
	}
	modern := a.Config().SDK >= 29
	if !modern {
		held, err := permission.Held(permission.WriteStorage)
		if err != nil {
			return nil, err
		}
		if !held {
			return nil, ErrNoPermission
		}
	}

	w := &Writer{}
	err := jni.Do(func(e *jni.Env) error {
		values, err := e.Make("android/content/ContentValues", jni.Sig(jni.TVoid))
		if err != nil {
			return err
		}
		put := func(key, value string) error { return putString(e, values, key, value) }
		if err := put(colName, name); err != nil {
			return err
		}
		if err := put(colMime, mime); err != nil {
			return err
		}
		if modern {
			path := kind.folder()
			if album != "" {
				path += "/" + album
			}
			if err := put(colPath, path); err != nil {
				return err
			}
			// Pending until it is written, so that nothing sees a half file.
			if err := putInt(e, values, colPending, 1); err != nil {
				return err
			}
		}

		class, field := kind.collection()
		collection, err := e.Constant(class, field, jni.TClass("android/net/Uri"))
		if err != nil {
			return err
		}
		resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver")))
		if err != nil {
			return err
		}
		uri, err := e.Invoke(resolver, "insert",
			jni.Sig(jni.TClass("android/net/Uri"), jni.TClass("android/net/Uri"),
				jni.TClass("android/content/ContentValues")),
			jni.Ref(collection), jni.Ref(values))
		if err != nil {
			return err
		}
		if uri.IsNil() {
			return fmt.Errorf("gallery: the library would not take a %s called %q",
				kind, name)
		}
		w.uri = e.Global(uri)
		if w.uriText, err = e.InvokeString(uri, "toString", jni.Sig(jni.TString)); err != nil {
			return err
		}

		// A descriptor rather than an OutputStream. Writing through Java
		// would mean a Java byte array per chunk; a descriptor is a file Go
		// writes to directly, at the speed of the disk.
		mode, err := e.String("w")
		if err != nil {
			return err
		}
		pfd, err := e.Invoke(resolver, "openFileDescriptor",
			jni.Sig(jni.TClass("android/os/ParcelFileDescriptor"),
				jni.TClass("android/net/Uri"), jni.TString),
			jni.Ref(uri), jni.Ref(mode))
		if err != nil {
			return err
		}
		if pfd.IsNil() {
			return errors.New("gallery: the library gave nothing to write to")
		}
		// detachFd hands the descriptor over: the Java object stops owning
		// it and closing the file is now Go's job.
		fd, err := e.InvokeInt(pfd, "detachFd", jni.Sig(jni.TInt))
		if err != nil {
			return err
		}
		w.file = os.NewFile(uintptr(fd), name)
		if w.file == nil {
			return errors.New("gallery: the descriptor the library gave is not usable")
		}
		return nil
	})
	if err != nil {
		w.abandon()
		return nil, err
	}
	return w, nil
}

// Write puts bytes in.
func (w *Writer) Write(p []byte) (int, error) {
	if w.closed {
		return 0, errors.New("gallery: the entry is closed")
	}
	return w.file.Write(p)
}

// URI is the content:// address of what is being written. It is what to hand
// another app, and it is not a path: on a modern device there may be no path
// at all.
func (w *Writer) URI() string { return w.uriText }

// Close finishes the file and publishes it. Until it is called nothing else
// on the device can see the entry.
func (w *Writer) Close() error {
	w.once.Do(func() {
		w.closed = true
		if w.file != nil {
			w.err = w.file.Close()
		}
		a := app.Current()
		if a == nil || a.Config().SDK < 29 {
			// Below Android 10 there is no pending flag, and the entry was
			// visible from the moment it was inserted.
			w.forget()
			return
		}
		if err := jni.Do(func(e *jni.Env) error {
			values, err := e.Make("android/content/ContentValues", jni.Sig(jni.TVoid))
			if err != nil {
				return err
			}
			if err := putInt(e, values, colPending, 0); err != nil {
				return err
			}
			resolver, err := e.Invoke(app.Context(), "getContentResolver",
				jni.Sig(jni.TClass("android/content/ContentResolver")))
			if err != nil {
				return err
			}
			_, err = e.InvokeInt(resolver, "update",
				jni.Sig(jni.TInt, jni.TClass("android/net/Uri"),
					jni.TClass("android/content/ContentValues"), jni.TString,
					jni.TArray(jni.TString)),
				jni.Ref(w.uri), jni.Ref(values), jni.Ref(jni.Object{}), jni.Ref(jni.Object{}))
			return err
		}); err != nil && w.err == nil {
			w.err = err
		}
		w.forget()
	})
	return w.err
}

// abandon deletes an entry that was made and never written.
func (w *Writer) abandon() {
	if w.uri.IsNil() {
		return
	}
	jni.Do(func(e *jni.Env) error {
		if resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver"))); err == nil {
			e.InvokeInt(resolver, "delete",
				jni.Sig(jni.TInt, jni.TClass("android/net/Uri"), jni.TString,
					jni.TArray(jni.TString)),
				jni.Ref(w.uri), jni.Ref(jni.Object{}), jni.Ref(jni.Object{}))
		}
		e.DeleteGlobal(w.uri)
		return nil
	})
	w.uri = jni.Object{}
}

func (w *Writer) forget() {
	if w.uri.IsNil() {
		return
	}
	jni.Do(func(e *jni.Env) error { e.DeleteGlobal(w.uri); return nil })
	w.uri = jni.Object{}
}

// Save writes bytes into the device's library in one call, and gives back
// the content:// address.
func Save(kind Kind, name, mime, album string, data []byte) (string, error) {
	w, err := Create(kind, name, mime, album)
	if err != nil {
		return "", err
	}
	if _, err := w.Write(data); err != nil {
		w.abandon()
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return w.URI(), nil
}

// SaveFrom copies from a reader, for something too big to hold in memory.
func SaveFrom(kind Kind, name, mime, album string, r io.Reader) (string, error) {
	w, err := Create(kind, name, mime, album)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(w, r); err != nil {
		w.abandon()
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return w.URI(), nil
}

func putString(e *jni.Env, values jni.Object, key, value string) error {
	k, err := e.String(key)
	if err != nil {
		return err
	}
	v, err := e.String(value)
	if err != nil {
		return err
	}
	return e.InvokeVoid(values, "put",
		jni.Sig(jni.TVoid, jni.TString, jni.TString), jni.Ref(k), jni.Ref(v))
}

// putInt goes through Integer.valueOf because ContentValues.put takes boxed
// numbers — there is no overload for a plain int.
func putInt(e *jni.Env, values jni.Object, key string, value int32) error {
	k, err := e.String(key)
	if err != nil {
		return err
	}
	boxed, err := e.Static("java/lang/Integer", "valueOf",
		jni.Sig(jni.TClass("java/lang/Integer"), jni.TInt), jni.Int(value))
	if err != nil {
		return err
	}
	return e.InvokeVoid(values, "put",
		jni.Sig(jni.TVoid, jni.TString, jni.TClass("java/lang/Integer")),
		jni.Ref(k), jni.Ref(boxed))
}
