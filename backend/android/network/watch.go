//go:build android

package network

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Watch calls f whenever the connection changes — a network appearing, the
// one in use going away, or what it can do changing, which is how a
// connection going from metered to unmetered is noticed.
//
// f runs on a thread of the system's choosing, not the app's, and is given
// the whole state rather than the difference, because the difference is
// almost never what a caller wants. Call stop when done; a callback left
// registered outlives the app's interest in it and goes on being called.
func Watch(f func(State)) (stop func(), err error) {
	token, release := app.Listen(watchHandler(f))

	var callback jni.Object
	err = jni.Do(func(e *jni.Env) error {
		cb, err := registerCallback(e, token)
		if err != nil {
			return err
		}
		callback = cb
		return nil
	})
	if err != nil {
		release()
		return func() {}, err
	}
	return watchStop(callback, release), nil
}
