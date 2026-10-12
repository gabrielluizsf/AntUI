//go:build android

package gallery

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// fillWriter inserts the library entry and opens somewhere to write it.
func fillWriter(e *jni.Env, w *Writer, kind Kind, name, mime, album string, modern bool) error {
	values, err := buildValues(e, kind, name, mime, album, modern)
	if err != nil {
		return err
	}
	uri, resolver, err := insertEntry(e, kind, values)
	if err != nil {
		return err
	}
	if uri.IsNil() {
		return fmt.Errorf("gallery: the library would not take a %s called %q",
			kind, name)
	}
	w.uri = e.Global(uri)
	if w.uriText, err = e.InvokeString(uri, "toString", jni.Sig(jni.TString)); err != nil {
		return err
	}
	file, err := openWriter(e, resolver, uri, name)
	if err != nil {
		return err
	}
	w.file = file
	return nil
}
