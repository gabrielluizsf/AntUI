//go:build android

package audio

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// registerNoisy registers the receiver for the becoming-noisy broadcast and
// returns the global reference to it.
func registerNoisy(e *jni.Env, token int64) (jni.Object, error) {
	r, err := app.Listener(e, "dev/antui/AntuiReceiver", token)
	if err != nil {
		return jni.Object{}, err
	}
	receiver := e.Global(r)

	action, err := e.ConstantString("android/media/AudioManager",
		"ACTION_AUDIO_BECOMING_NOISY")
	if err != nil {
		return jni.Object{}, err
	}
	jaction, err := e.String(action)
	if err != nil {
		return jni.Object{}, err
	}
	filter, err := e.Make("android/content/IntentFilter",
		jni.Sig(jni.TVoid, jni.TString), jni.Ref(jaction))
	if err != nil {
		return jni.Object{}, err
	}
	_, err = e.Invoke(app.Context(), "registerReceiver",
		jni.Sig(jni.TClass("android/content/Intent"),
			jni.TClass("android/content/BroadcastReceiver"),
			jni.TClass("android/content/IntentFilter")),
		jni.Ref(receiver), jni.Ref(filter))
	if err != nil {
		return jni.Object{}, err
	}
	return receiver, nil
}
