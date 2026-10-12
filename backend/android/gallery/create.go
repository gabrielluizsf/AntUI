//go:build android

package gallery

import (
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
		return fillWriter(e, w, kind, name, mime, album, modern)
	})
	if err != nil {
		w.abandon()
		return nil, err
	}
	return w, nil
}
