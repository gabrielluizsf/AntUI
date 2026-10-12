//go:build android

package intent

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// ShareFile offers a file to any app that takes its kind.
//
// The address has to be a content:// one from a provider — see
// [antui/backend/android/share], which is what makes one. A path will not do, and
// handing one over throws rather than failing quietly, which is the platform
// being firm about something that used to be a way to leak private files.
func ShareFile(uri, mime, title string) error {
	return jni.Do(func(e *jni.Env) error {
		in, err := shareIntent(e, uri, mime)
		if err != nil {
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
