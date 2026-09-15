//go:build android

package audio

import (
	"errors"
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Kind is how long the app expects to want the sound for. Asking for the
// right one is not politeness: it is what tells the music app whether to
// stop or to duck, and getting it wrong is why some apps kill a podcast to
// play a notification.
type Kind int32

const (
	// Permanent is for something that plays until the user stops it — a
	// game, a music player. Whatever was playing before stops.
	Permanent Kind = 1
	// Brief is for something short that the other app should pause for and
	// resume after: a spoken direction, a video clip.
	Brief Kind = 2
	// BriefDucking is for something short that can play *over* the other
	// app, which turns itself down rather than stopping. It is what almost
	// every notification sound should ask for.
	BriefDucking Kind = 3
)

// Change is what happened to the app's hold on the sound.
type Change int32

const (
	// Gained means the app has it, or has it back.
	Gained Change = 1
	// Lost means something else took it for good. Stop, and do not expect
	// it back.
	Lost Change = -1
	// LostBriefly means something else took it for a moment — a call, a
	// spoken direction. Pause, and wait for Gained.
	LostBriefly Change = -2
	// Duck means something else is playing over the app. Turn down, do not
	// stop; the other sound is short and stopping would be worse.
	Duck Change = -3
)

// String names the change, which is what a log line wants.
func (c Change) String() string {
	switch c {
	case Gained:
		return "gained"
	case Lost:
		return "lost"
	case LostBriefly:
		return "lost briefly"
	case Duck:
		return "duck"
	}
	return "unknown"
}

// streamMusic is AudioManager.STREAM_MUSIC, which is the one everything that
// is not a phone call plays on.
const streamMusic int32 = 3

// focusGranted is AUDIOFOCUS_REQUEST_GRANTED.
const focusGranted int32 = 1

// ErrRefused is the system saying no — a phone call in progress, or another
// app holding the sound in a way that cannot be interrupted.
var ErrRefused = errors.New("audio: the system refused audio focus")

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

// OnHeadsetUnplugged calls f when the headphones come out or the bluetooth
// speaker goes away.
//
// The rule the platform expects, and the one users notice when an app breaks
// it: **pause**. Sound that was private a moment ago is about to come out of
// the speaker in a room with other people in it. It is not a suggestion —
// the broadcast exists for nothing else.
func OnHeadsetUnplugged(f func()) (stop func(), err error) {
	token, releaseToken := app.Listen(func(e *jni.Env, method string, args []jni.Object) {
		if method == "onReceive" && f != nil {
			f()
		}
	})

	var receiver jni.Object
	err = jni.Do(func(e *jni.Env) error {
		r, err := app.Listener(e, "dev/antui/AntuiReceiver", token)
		if err != nil {
			return err
		}
		receiver = e.Global(r)

		action, err := e.ConstantString("android/media/AudioManager",
			"ACTION_AUDIO_BECOMING_NOISY")
		if err != nil {
			return err
		}
		jaction, err := e.String(action)
		if err != nil {
			return err
		}
		filter, err := e.Make("android/content/IntentFilter",
			jni.Sig(jni.TVoid, jni.TString), jni.Ref(jaction))
		if err != nil {
			return err
		}
		_, err = e.Invoke(app.Context(), "registerReceiver",
			jni.Sig(jni.TClass("android/content/Intent"),
				jni.TClass("android/content/BroadcastReceiver"),
				jni.TClass("android/content/IntentFilter")),
			jni.Ref(receiver), jni.Ref(filter))
		return err
	})
	if err != nil {
		releaseToken()
		return nil, err
	}

	var once sync.Once
	return func() {
		once.Do(func() {
			jni.Do(func(e *jni.Env) error {
				e.InvokeVoid(app.Context(), "unregisterReceiver",
					jni.Sig(jni.TVoid, jni.TClass("android/content/BroadcastReceiver")),
					jni.Ref(receiver))
				e.DeleteGlobal(receiver)
				return nil
			})
			releaseToken()
		})
	}, nil
}
