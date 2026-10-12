//go:build android

package speech

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// makeEngine builds the TextToSpeech object and wires the progress listener,
// on the UI thread.
func makeEngine(e *jni.Env, s *Speaker, engine string, initToken, progressToken int64) error {
	listener, err := app.Proxy(e, initToken,
		"android.speech.tts.TextToSpeech$OnInitListener")
	if err != nil {
		return err
	}
	tts, err := newTextToSpeech(e, engine, listener)
	if err != nil {
		return err
	}
	s.tts = e.Global(tts)

	progress, err := app.Listener(e, "dev/antui/AntuiUtterance", progressToken)
	if err != nil {
		return err
	}
	// **InvokeInt and not InvokeVoid.** This one returns an int — SUCCESS or
	// ERROR — and calling a method that returns a value as though it
	// returned none is undefined behaviour: the VM's CallVoidMethodA reads
	// the frame back differently. It works until it does not, and nothing
	// says so unless CheckJNI is on, which it is not by default. The
	// signature here was right all along; the call was not.
	_, err = e.InvokeInt(s.tts, "setOnUtteranceProgressListener",
		jni.Sig(jni.TInt,
			jni.TClass("android/speech/tts/UtteranceProgressListener")),
		jni.Ref(progress))
	return err
}
