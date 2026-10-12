//go:build android

package picker

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Forget gives a kept permission back.
func Forget(uri string, write bool) error {
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
		return e.InvokeVoid(resolver, "releasePersistableUriPermission",
			jni.Sig(jni.TVoid, jni.TClass("android/net/Uri"), jni.TInt),
			jni.Ref(parsed), jni.Int(flags))
	})
}
