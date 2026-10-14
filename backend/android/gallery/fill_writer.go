//go:build android

package gallery

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// fillWriter inserts the library entry and opens somewhere to write it.
func fillWriter(e *jni.Env, w *Writer, entry Entry, modern bool) error {
	values, err := buildValues(e, entry, modern)
	if err != nil {
		return err
	}
	uri, resolver, err := insertEntry(e, entry.Kind, values)
	if err != nil {
		return err
	}
	if uri.IsNil() {
		return fmt.Errorf("gallery: the library would not take a %s called %q",
			entry.Kind, entry.Name)
	}
	w.uri = e.Global(uri)
	if w.uriText, err = e.InvokeString(uri, "toString", jni.Sig(jni.TString)); err != nil {
		return err
	}
	file, err := openWriter(e, resolver, uri, entry.Name)
	if err != nil {
		return err
	}
	w.file = file
	return nil
}
