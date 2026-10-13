//go:build android

package app

import "github.com/gabrielluizsf/antui/backend/android/jni"

// Proxy makes a Java object that implements the named interfaces and sends
// every call to the handler registered under token.
//
// The names are dotted, as Class.forName takes them:
//
//	app.Proxy(e, token, "android.content.DialogInterface$OnClickListener")
//
// It only works for **interfaces**. An abstract class cannot be proxied, and
// the two Android insists on — BroadcastReceiver and NetworkCallback — are
// written out in the shim instead. The object comes back as a local
// reference; anything that outlives the call needs [jni.Env.Global].
func Proxy(e *jni.Env, token int64, interfaces ...string) (jni.Object, error) {
	names, err := e.NewStrings(interfaces)
	if err != nil {
		return jni.Object{}, err
	}
	return e.Static("dev/antui/AntuiProxy", "create",
		jni.Sig(jni.TObject, jni.TLong, jni.TArray(jni.TString)),
		jni.Long(token), jni.Ref(names))
}
