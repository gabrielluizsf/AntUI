//go:build android

package gallery

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Close finishes the file and publishes it. Until it is called nothing else
// on the device can see the entry.
func (w *Writer) Close() error {
	w.once.Do(func() {
		w.closed = true
		if w.file != nil {
			w.err = w.file.Close()
		}
		a := app.Current()
		if a == nil || a.Config().SDK < 29 {
			// Below Android 10 there is no pending flag, and the entry was
			// visible from the moment it was inserted.
			w.forget()
			return
		}
		if err := jni.Do(func(e *jni.Env) error {
			values, err := e.Make("android/content/ContentValues", jni.Sig(jni.TVoid))
			if err != nil {
				return err
			}
			if err := putInt(e, values, colPending, 0); err != nil {
				return err
			}
			resolver, err := e.Invoke(app.Context(), "getContentResolver",
				jni.Sig(jni.TClass("android/content/ContentResolver")))
			if err != nil {
				return err
			}
			_, err = e.InvokeInt(resolver, "update",
				jni.Sig(jni.TInt, jni.TClass("android/net/Uri"),
					jni.TClass("android/content/ContentValues"), jni.TString,
					jni.TArray(jni.TString)),
				jni.Ref(w.uri), jni.Ref(values), jni.Ref(jni.Object{}), jni.Ref(jni.Object{}))
			return err
		}); err != nil && w.err == nil {
			w.err = err
		}
		w.forget()
	})
	return w.err
}
