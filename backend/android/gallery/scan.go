//go:build android

package gallery

import (
	"errors"
	"time"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Scan tells the device to look at a file that is already on disk and put it
// in the library.
//
// It is the old way, and it is here because a file an app wrote itself —
// into its own external directory, say — is invisible to the gallery until
// something says it is there. On Android 10 and up prefer [Create]: an entry
// made through the library is owned, published and cleaned up properly,
// where a scanned file is a file the app happens to have left somewhere.
//
// It waits for the answer, which arrives as the content:// address the file
// was given, or empty when the scanner would not take it — a format it does
// not know, or a path it may not read.
func Scan(path, mime string) (string, error) {
	answer := make(chan string, 1)
	token, release := app.Listen(func(e *jni.Env, method string, args []jni.Object) {
		if method != "onScanCompleted" {
			return
		}
		out := ""
		if len(args) > 1 && !args[1].IsNil() {
			out, _ = e.InvokeString(args[1], "toString", jni.Sig(jni.TString))
		}
		select {
		case answer <- out:
		default:
		}
	})
	defer release()

	err := jni.Do(func(e *jni.Env) error {
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
	})
	if err != nil {
		return "", err
	}
	select {
	case uri := <-answer:
		return uri, nil
	case <-time.After(15 * time.Second):
		return "", errors.New("gallery: the scanner never answered")
	}
}
