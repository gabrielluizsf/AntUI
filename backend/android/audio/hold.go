//go:build android

package audio

import (
	"sync"

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
	token, releaseToken := app.Listen(func(e *jni.Env, method string, args []jni.Object) {
		if method != "onAudioFocusChange" || len(args) == 0 {
			return
		}
		which, err := e.InvokeInt(args[0], "intValue", jni.Sig(jni.TInt))
		if err != nil {
			return
		}
		if onChange != nil {
			onChange(Change(which))
		}
	})

	var listener jni.Object
	err = jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, "audio")
		if err != nil {
			return err
		}
		// OnAudioFocusChangeListener is an interface, so a proxy is enough.
		l, err := app.Proxy(e, token,
			"android.media.AudioManager$OnAudioFocusChangeListener")
		if err != nil {
			return err
		}
		// Global: the platform holds it, and abandoning focus needs the very
		// same object — a different one that behaves the same does not count.
		listener = e.Global(l)

		// The request that takes a listener directly is deprecated in favour
		// of AudioFocusRequest, which is API 26. This one works on every
		// version this library supports, and does the same thing.
		got, err := e.InvokeInt(manager, "requestAudioFocus",
			jni.Sig(jni.TInt,
				jni.TClass("android/media/AudioManager$OnAudioFocusChangeListener"),
				jni.TInt, jni.TInt),
			jni.Ref(listener), jni.Int(streamMusic), jni.Int(int32(kind)))
		if err != nil {
			return err
		}
		if got != focusGranted {
			return ErrRefused
		}
		return nil
	})
	if err != nil {
		releaseToken()
		return nil, err
	}

	var once sync.Once
	return func() {
		once.Do(func() {
			jni.Do(func(e *jni.Env) error {
				if manager, err := app.Service(e, "audio"); err == nil {
					e.InvokeInt(manager, "abandonAudioFocus",
						jni.Sig(jni.TInt,
							jni.TClass("android/media/AudioManager$OnAudioFocusChangeListener")),
						jni.Ref(listener))
				}
				e.DeleteGlobal(listener)
				return nil
			})
			releaseToken()
		})
	}, nil
}
