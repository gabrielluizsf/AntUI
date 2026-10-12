//go:build android

package audio

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// OnHeadsetUnplugged calls f when the headphones come out or the bluetooth
// speaker goes away.
//
// The rule the platform expects, and the one users notice when an app breaks
// it: **pause**. Sound that was private a moment ago is about to come out of
// the speaker in a room with other people in it. It is not a suggestion —
// the broadcast exists for nothing else.
func OnHeadsetUnplugged(f func()) (stop func(), err error) {
	token, releaseToken := app.Listen(unpluggedHandler(f))

	var receiver jni.Object
	err = jni.Do(func(e *jni.Env) error {
		r, err := registerNoisy(e, token)
		if err != nil {
			return err
		}
		receiver = r
		return nil
	})
	if err != nil {
		releaseToken()
		return nil, err
	}
	return unpluggedStop(receiver, releaseToken), nil
}
