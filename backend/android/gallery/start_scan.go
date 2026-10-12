//go:build android

package gallery

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// startScan asks MediaScannerConnection to scan the file, with the listener
// that will be told the answer.
func startScan(e *jni.Env, token int64, path, mime string) error {
	paths, err := e.NewStrings([]string{path})
	if err != nil {
		return err
	}
	mimes, err := e.NewStrings([]string{mime})
	if err != nil {
		return err
	}
	listener, err := app.Proxy(e, token,
		"android.media.MediaScannerConnection$OnScanCompletedListener")
	if err != nil {
		return err
	}
	return e.StaticVoid("android/media/MediaScannerConnection", "scanFile",
		jni.Sig(jni.TVoid, jni.TContext, jni.TArray(jni.TString),
			jni.TArray(jni.TString),
			jni.TClass("android/media/MediaScannerConnection$OnScanCompletedListener")),
		jni.Ref(app.Context()), jni.Ref(paths), jni.Ref(mimes), jni.Ref(listener))
}
