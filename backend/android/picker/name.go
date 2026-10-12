//go:build android

package picker

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Name is what the user would call the file — the name shown in the picker,
// which is not in the address and has to be asked for.
func Name(uri string) (string, error) {
	var out string
	err := jni.Do(func(e *jni.Env) error {
		resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver")))
		if err != nil {
			return err
		}
		parsed, err := parse(e, uri)
		if err != nil {
			return err
		}
		cursor, err := e.Invoke(resolver, "query",
			jni.Sig(jni.TClass("android/database/Cursor"),
				jni.TClass("android/net/Uri"), jni.TArray(jni.TString),
				jni.TString, jni.TArray(jni.TString), jni.TString),
			jni.Ref(parsed), jni.Ref(jni.Object{}), jni.Ref(jni.Object{}),
			jni.Ref(jni.Object{}), jni.Ref(jni.Object{}))
		if err != nil || cursor.IsNil() {
			return err
		}
		defer e.InvokeVoid(cursor, "close", jni.Sig(jni.TVoid))

		moved, err := e.InvokeBool(cursor, "moveToFirst", jni.Sig(jni.TBool))
		if err != nil || !moved {
			return err
		}
		col, err := e.String("_display_name")
		if err != nil {
			return err
		}
		at, err := e.InvokeInt(cursor, "getColumnIndex",
			jni.Sig(jni.TInt, jni.TString), jni.Ref(col))
		if err != nil || at < 0 {
			return err
		}
		out, err = e.InvokeString(cursor, "getString",
			jni.Sig(jni.TString, jni.TInt), jni.Int(at))
		return err
	})
	return out, err
}
