//go:build android

package picker

import "github.com/gabrielluizsf/antui/backend/android/jni"

// build makes the intent.
//
// ACTION_OPEN_DOCUMENT rather than ACTION_GET_CONTENT: the first goes to the
// system's own picker and gives an address that keeps working, and the second
// goes to whichever app claims the type and gives one that may not outlive
// the call. The system picker is also the one that shows cloud accounts.
func build(mime string, many bool) (jni.Object, error) {
	var out jni.Object
	err := jni.Do(func(e *jni.Env) error {
		in, err := openDocumentIntent(e, mime, many)
		if err != nil {
			return err
		}
		out = e.Global(in)
		return nil
	})
	return out, err
}
