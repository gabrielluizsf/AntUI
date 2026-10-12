//go:build android

package gallery

import (
	"errors"
	"fmt"
	"os"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/permission"
)

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
