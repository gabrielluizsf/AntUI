//go:build android

package gallery

import (
	"errors"
	"os"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// openWriter returns a file the entry can be written to.
//
// A descriptor rather than an OutputStream. Writing through Java would mean
// a Java byte array per chunk; a descriptor is a file Go writes to directly,
// at the speed of the disk.
func openWriter(e *jni.Env, resolver, uri jni.Object, name string) (*os.File, error) {
	mode, err := e.String("w")
	if err != nil {
		return nil, err
	}
	pfd, err := e.Invoke(resolver, "openFileDescriptor",
		jni.Sig(jni.TClass("android/os/ParcelFileDescriptor"),
			jni.TClass("android/net/Uri"), jni.TString),
		jni.Ref(uri), jni.Ref(mode))
	if err != nil {
		return nil, err
	}
	if pfd.IsNil() {
		return nil, errors.New("gallery: the library gave nothing to write to")
	}
	// detachFd hands the descriptor over: the Java object stops owning it
	// and closing the file is now Go's job.
	fd, err := e.InvokeInt(pfd, "detachFd", jni.Sig(jni.TInt))
	if err != nil {
		return nil, err
	}
	file := os.NewFile(uintptr(fd), name)
	if file == nil {
		return nil, errors.New("gallery: the descriptor the library gave is not usable")
	}
	return file, nil
}
