//go:build android

package gallery

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// abandon deletes an entry that was made and never written.
func (w *Writer) abandon() {
	if w.uri.IsNil() {
		return
	}
	jni.Do(func(e *jni.Env) error {
		if resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver"))); err == nil {
			e.InvokeInt(resolver, "delete",
				jni.Sig(jni.TInt, jni.TClass("android/net/Uri"), jni.TString,
					jni.TArray(jni.TString)),
				jni.Ref(w.uri), jni.Ref(jni.Object{}), jni.Ref(jni.Object{}))
		}
		e.DeleteGlobal(w.uri)
		return nil
	})
	w.uri = jni.Object{}
}
