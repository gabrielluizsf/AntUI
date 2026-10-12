//go:build android

package audio

import (
	"sync"

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
