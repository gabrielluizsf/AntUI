//go:build android

package intent

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// shareIntent builds the ACTION_SEND intent that carries one content://
// address, with the temporary grant the receiving app needs to read it.
func shareIntent(e *jni.Env, uri, mime string) (jni.Object, error) {
	action, err := e.String(ActionSend)
	if err != nil {
		return jni.Object{}, err
	}
	in, err := e.Make("android/content/Intent",
		jni.Sig(jni.TVoid, jni.TString), jni.Ref(action))
	if err != nil {
		return jni.Object{}, err
	}
	jmime, err := e.String(mime)
	if err != nil {
		return jni.Object{}, err
	}
	if _, err := e.Invoke(in, "setType",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString),
		jni.Ref(jmime)); err != nil {
		return jni.Object{}, err
	}
	juri, err := e.String(uri)
	if err != nil {
		return jni.Object{}, err
	}
	parsed, err := e.Static("android/net/Uri", "parse",
		jni.Sig(jni.TClass("android/net/Uri"), jni.TString), jni.Ref(juri))
	if err != nil {
		return jni.Object{}, err
	}
	// EXTRA_STREAM carries a Uri and not a string. Putting the text of one
	// in works as far as the compiler is concerned and produces a share
	// every app refuses.
	key, err := e.String("android.intent.extra.STREAM")
	if err != nil {
		return jni.Object{}, err
	}
	if _, err := e.Invoke(in, "putExtra",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString,
			jni.TClass("android/os/Parcelable")),
		jni.Ref(key), jni.Ref(parsed)); err != nil {
		return jni.Object{}, err
	}
	if _, err := e.Invoke(in, "addFlags",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TInt),
		jni.Int(grantRead)); err != nil {
		return jni.Object{}, err
	}
	return in, nil
}
