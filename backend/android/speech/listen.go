//go:build android

package speech

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/permission"
)

// Listen starts the recogniser and reports what it hears.
//
// f runs on the UI thread, which is where the engine answers. Every call has
// either a result or an error, never both.
//
// It needs RECORD_AUDIO. On almost every device the sound goes to a server to
// be recognised, so an app that listens is an app that sends what was said
// somewhere — which the user should be told, in the app and not only in a
// privacy policy.
func Listen(opt Options, f func(Heard, error)) (stop func(), err error) {
	held, err := permission.Held(permission.Microphone)
	if err != nil {
		return nil, err
	}
	if !held {
		return nil, ErrNoPermission
	}
	if !Available() {
		return nil, ErrNoRecogniser
	}

	l := &listenState{opt: opt, report: f}
	token, release := app.Listen(l.onEvent)

	var setupErr error
	if e := app.RunOnUISync(func() {
		setupErr = jni.Do(l.setup(token))
	}); e != nil {
		release()
		return nil, e
	}
	if setupErr != nil {
		release()
		return nil, setupErr
	}
	return l.stop(release), nil
}
