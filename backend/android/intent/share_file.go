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
		action, err := e.String(ActionSend)
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
		juri, err := e.String(uri)
		if err != nil {
			return err
		}
		parsed, err := e.Static("android/net/Uri", "parse",
			jni.Sig(jni.TClass("android/net/Uri"), jni.TString), jni.Ref(juri))
		if err != nil {
			return err
		}
		// EXTRA_STREAM carries a Uri and not a string. Putting the text of
		// one in works as far as the compiler is concerned and produces a
		// share every app refuses.
		key, err := e.String("android.intent.extra.STREAM")
		if err != nil {
			return err
		}
		if _, err := e.Invoke(in, "putExtra",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString,
				jni.TClass("android/os/Parcelable")),
			jni.Ref(key), jni.Ref(parsed)); err != nil {
			return err
		}
		if _, err := e.Invoke(in, "addFlags",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TInt),
			jni.Int(grantRead)); err != nil {
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
