//go:build android

package picker

import (
	"errors"
	"os"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

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
