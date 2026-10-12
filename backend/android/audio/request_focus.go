//go:build android

package audio

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// requestFocus asks for the focus and returns the global listener the
// platform now holds.
func requestFocus(e *jni.Env, token int64, kind Kind) (jni.Object, error) {
	manager, err := app.Service(e, "audio")
	if err != nil {
		return jni.Object{}, err
	}
	// OnAudioFocusChangeListener is an interface, so a proxy is enough.
	l, err := app.Proxy(e, token,
		"android.media.AudioManager$OnAudioFocusChangeListener")
	if err != nil {
		return jni.Object{}, err
	}
	// Global: the platform holds it, and abandoning focus needs the very
	// same object — a different one that behaves the same does not count.
	listener := e.Global(l)

	// The request that takes a listener directly is deprecated in favour of
	// AudioFocusRequest, which is API 26. This one works on every version
	// this library supports, and does the same thing.
	got, err := e.InvokeInt(manager, "requestAudioFocus",
		jni.Sig(jni.TInt,
			jni.TClass("android/media/AudioManager$OnAudioFocusChangeListener"),
			jni.TInt, jni.TInt),
		jni.Ref(listener), jni.Int(streamMusic), jni.Int(int32(kind)))
	if err != nil {
		return jni.Object{}, err
	}
	if got != focusGranted {
		return jni.Object{}, ErrRefused
	}
	return listener, nil
}
