//go:build android

package speech

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// setup creates the recogniser and starts it. It runs on the UI thread.
func (l *listenState) setup(token int64) func(*jni.Env) error {
	return func(e *jni.Env) error {
		c, err := e.Class("android/speech/SpeechRecognizer")
		if err != nil {
			return err
		}
		create, err := e.StaticMethod(c, "createSpeechRecognizer",
			jni.Sig(jni.TClass("android/speech/SpeechRecognizer"), jni.TContext))
		if err != nil {
			return err
		}
		r, err := e.CallStaticObject(c, create, jni.Ref(app.Context()))
		if err != nil {
			return err
		}
		if r.IsNil() {
			return ErrNoRecogniser
		}
		l.recogniser = e.Global(r)

		listener, err := app.Proxy(e, token, "android.speech.RecognitionListener")
		if err != nil {
			return err
		}
		if err := e.InvokeVoid(l.recogniser, "setRecognitionListener",
			jni.Sig(jni.TVoid, jni.TClass("android/speech/RecognitionListener")),
			jni.Ref(listener)); err != nil {
			return err
		}

		in, err := buildIntent(e, l.opt)
		if err != nil {
			return err
		}
		l.intent = e.Global(in)
		return startListening(e, l.recogniser, l.intent)
	}
}
