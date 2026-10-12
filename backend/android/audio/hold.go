//go:build android

package audio

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Hold asks the system for the right to make sound, and reports what happens
// to it afterwards.
//
// An app that plays anything at all should ask. An app that does not will be
// heard over the music the user was already listening to, and will go on
// being heard through a phone call.
//
// onChange runs on the UI thread. It is where an app pauses, ducks and
// resumes; ignoring it is the same as not asking in the first place.
func Hold(kind Kind, onChange func(Change)) (release func(), err error) {
	token, releaseToken := app.Listen(focusListener(onChange))

	var listener jni.Object
	err = jni.Do(func(e *jni.Env) error {
		l, err := requestFocus(e, token, kind)
		if err != nil {
			return err
		}
		listener = l
		return nil
	})
	if err != nil {
		releaseToken()
		return nil, err
	}
	return focusRelease(listener, releaseToken), nil
}
