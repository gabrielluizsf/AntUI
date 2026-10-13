//go:build android

package gallery

import "github.com/gabrielluizsf/antui/backend/android/jni"

// putInt goes through Integer.valueOf because ContentValues.put takes boxed
// numbers — there is no overload for a plain int.
func putInt(e *jni.Env, values jni.Object, key string, value int32) error {
	k, err := e.String(key)
	if err != nil {
		return err
	}
	boxed, err := e.Static("java/lang/Integer", "valueOf",
		jni.Sig(jni.TClass("java/lang/Integer"), jni.TInt), jni.Int(value))
	if err != nil {
		return err
	}
	return e.InvokeVoid(values, "put",
		jni.Sig(jni.TVoid, jni.TString, jni.TClass("java/lang/Integer")),
		jni.Ref(k), jni.Ref(boxed))
}
