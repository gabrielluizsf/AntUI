//go:build android

package intent

import "github.com/gabrielluizsf/antui/backend/android/jni"

// View opens something with the app that handles its type. mime may be empty,
// in which case the system guesses from the address.
func View(uri, mime string) error {
	return jni.Do(func(e *jni.Env) error {
		in, err := viewIntent(e, uri, mime)
		if err != nil {
			return err
		}
		return start(e, in)
	})
}
