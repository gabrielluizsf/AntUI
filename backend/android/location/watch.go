//go:build android

package location

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Watch reports where the device is, as it moves.
//
// f runs on the UI thread, because that is the thread the listener was
// registered from and the platform calls back on the caller's looper. Do not
// take long in it.
//
// **This costs battery**, and the two providers cost very differently: the
// satellites can take a few per cent an hour and the towers almost nothing.
// Stop it the moment the app stops caring.
func Watch(opt Options, f func(Fix)) (stop func(), err error) {
	ok, err := held()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNoPermission
	}
	providers := []string{opt.Provider}
	if opt.Provider == "" {
		providers = []string{GPS, Network}
	}

	token, release := app.Listen(watchListener(f))

	var listener jni.Object
	var setupErr error
	if e := app.RunOnUISync(func() {
		setupErr = jni.Do(func(e *jni.Env) error {
			l, err := requestUpdates(e, token, opt, providers)
			if err != nil {
				return err
			}
			listener = l
			return nil
		})
	}); e != nil {
		release()
		return nil, e
	}
	if setupErr != nil {
		release()
		return nil, setupErr
	}
	return watchStop(listener, release), nil
}
