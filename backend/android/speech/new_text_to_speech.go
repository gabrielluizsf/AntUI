//go:build android

package speech

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// newTextToSpeech builds the engine, with a named engine when one was asked
// for and the device default otherwise.
func newTextToSpeech(e *jni.Env, engine string, listener jni.Object) (jni.Object, error) {
	if engine == "" {
		return e.Make("android/speech/tts/TextToSpeech",
			jni.Sig(jni.TVoid, jni.TContext,
				jni.TClass("android/speech/tts/TextToSpeech$OnInitListener")),
			jni.Ref(app.Context()), jni.Ref(listener))
	}
	jengine, err := e.String(engine)
	if err != nil {
		return jni.Object{}, err
	}
	return e.Make("android/speech/tts/TextToSpeech",
		jni.Sig(jni.TVoid, jni.TContext,
			jni.TClass("android/speech/tts/TextToSpeech$OnInitListener"),
			jni.TString),
		jni.Ref(app.Context()), jni.Ref(listener), jni.Ref(jengine))
}
