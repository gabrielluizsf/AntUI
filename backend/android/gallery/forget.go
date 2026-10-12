//go:build android

package gallery

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func (w *Writer) forget() {
	if w.uri.IsNil() {
		return
	}
	jni.Do(func(e *jni.Env) error { e.DeleteGlobal(w.uri); return nil })
	w.uri = jni.Object{}
}
