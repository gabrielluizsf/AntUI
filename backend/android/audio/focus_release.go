//go:build android

package audio

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// focusRelease returns the function that abandons focus, safe to call more
// than once.
func focusRelease(listener jni.Object, releaseToken func()) func() {
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
	}
}
