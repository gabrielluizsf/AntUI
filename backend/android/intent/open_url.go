//go:build android

package intent

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// OpenURL hands a link to whatever the user opens links with — a browser, or
// the app that claims that address.
func OpenURL(url string) error {
	return jni.Do(func(e *jni.Env) error {
		in, err := viewIntent(e, url, "")
		if err != nil {
			return err
		}
		return start(e, in)
	})
}
