//go:build android

package picker

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Create asks the user where to put a new file, and gives back its address.
//
// name is what the picker suggests; the user may change it, and what comes
// back is whatever they settled on. Nothing is written — the file exists and
// is empty, and [Writer] is how it is filled.
func Create(name, mime string) (string, error) {
	var in jni.Object
	err := jni.Do(func(e *jni.Env) error {
		var err error
		in, err = document(e, "android.intent.action.CREATE_DOCUMENT", mime, name, false)
		return err
	})
	if err != nil {
		return "", err
	}
	return single(in)
}
