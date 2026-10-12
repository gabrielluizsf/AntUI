//go:build android

package picker

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Keep asks the system to remember that the user chose this address, so that
// it still works after the app has been closed and opened again.
//
// Without it a chosen file works until the process ends and then stops,
// which looks to a user exactly like the app losing their document. There is
// a limit on how many an app may keep — a few hundred — and the oldest are
// dropped, so keeping every file a user ever opened is not a plan.
func Keep(uri string, write bool) error {
	flags := grantRead
	if write {
		flags |= grantWrite
	}
	return jni.Do(func(e *jni.Env) error {
		resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver")))
		if err != nil {
			return err
		}
		parsed, err := parse(e, uri)
		if err != nil {
			return err
		}
		return e.InvokeVoid(resolver, "takePersistableUriPermission",
			jni.Sig(jni.TVoid, jni.TClass("android/net/Uri"), jni.TInt),
			jni.Ref(parsed), jni.Int(flags))
	})
}
