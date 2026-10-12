//go:build android

package gallery

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// queryCursor runs the query and returns the cursor, with the collection it
// was made from — readItem needs it to build each item's address.
func queryCursor(e *jni.Env, kind Kind, projection, where, args, order jni.Object) (jni.Object, jni.Object, error) {
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
	cursor, err := e.Invoke(resolver, "query",
		jni.Sig(jni.TClass("android/database/Cursor"),
			jni.TClass("android/net/Uri"), jni.TArray(jni.TString),
			jni.TString, jni.TArray(jni.TString), jni.TString),
		jni.Ref(collection), jni.Ref(projection), jni.Ref(where),
		jni.Ref(args), jni.Ref(order))
	if err != nil {
		return jni.Object{}, jni.Object{}, err
	}
	return cursor, collection, nil
}
