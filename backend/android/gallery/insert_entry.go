//go:build android

package gallery

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// insertEntry inserts the row and returns the content:// address it was
// given, with the resolver it was inserted through.
func insertEntry(e *jni.Env, kind Kind, values jni.Object) (jni.Object, jni.Object, error) {
	class, field := kind.collection()
	collection, err := e.Constant(class, field, jni.TClass("android/net/Uri"))
	if err != nil {
		return jni.Object{}, jni.Object{}, err
	}
	resolver, err := e.Invoke(app.Context(), "getContentResolver",
		jni.Sig(jni.TClass("android/content/ContentResolver")))
	if err != nil {
		return jni.Object{}, jni.Object{}, err
	}
	uri, err := e.Invoke(resolver, "insert",
		jni.Sig(jni.TClass("android/net/Uri"), jni.TClass("android/net/Uri"),
			jni.TClass("android/content/ContentValues")),
		jni.Ref(collection), jni.Ref(values))
	if err != nil {
		return jni.Object{}, jni.Object{}, err
	}
	return uri, resolver, nil
}
