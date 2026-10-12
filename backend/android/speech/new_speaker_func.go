//go:build android

package speech

import (
	"errors"
	"time"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func newSpeaker(engine string) (*Speaker, error) {
	s := &Speaker{waits: map[string]chan error{}}

	ready := make(chan int32, 1)
	initToken, releaseInit := app.Listen(func(e *jni.Env, method string, args []jni.Object) {
		if method != "onInit" || len(args) == 0 {
			return
		}
		status, err := e.InvokeInt(args[0], "intValue", jni.Sig(jni.TInt))
		if err != nil {
			status = -1
		}
		select {
		case ready <- status:
		default:
		}
	})
	s.release = append(s.release, releaseInit)

	progressToken, releaseProgress := app.Listen(s.onProgress)
	s.release = append(s.release, releaseProgress)

	var err error
	if e := app.RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			listener, err := app.Proxy(e, initToken,
				"android.speech.tts.TextToSpeech$OnInitListener")
			if err != nil {
				return err
			}
			var tts jni.Object
			if engine == "" {
				tts, err = e.Make("android/speech/tts/TextToSpeech",
					jni.Sig(jni.TVoid, jni.TContext,
						jni.TClass("android/speech/tts/TextToSpeech$OnInitListener")),
					jni.Ref(app.Context()), jni.Ref(listener))
			} else {
				jengine, err2 := e.String(engine)
				if err2 != nil {
					return err2
				}
				tts, err = e.Make("android/speech/tts/TextToSpeech",
					jni.Sig(jni.TVoid, jni.TContext,
						jni.TClass("android/speech/tts/TextToSpeech$OnInitListener"),
						jni.TString),
					jni.Ref(app.Context()), jni.Ref(listener), jni.Ref(jengine))
			}
			if err != nil {
				return err
			}
			s.tts = e.Global(tts)

			progress, err := app.Listener(e, "dev/antui/AntuiUtterance", progressToken)
			if err != nil {
				return err
			}
			// **InvokeInt and not InvokeVoid.** This one returns an int —
			// SUCCESS or ERROR — and calling a method that returns a value
			// as though it returned none is undefined behaviour: the VM's
			// CallVoidMethodA reads the frame back differently. It works
			// until it does not, and nothing says so unless CheckJNI is on,
			// which it is not by default. The signature here was right all
			// along; the call was not.
			_, err = e.InvokeInt(s.tts, "setOnUtteranceProgressListener",
				jni.Sig(jni.TInt,
					jni.TClass("android/speech/tts/UtteranceProgressListener")),
				jni.Ref(progress))
			return err
		})
	}); e != nil {
		s.Close()
		return nil, e
	}
	if err != nil {
		s.Close()
		return nil, err
	}

	select {
	case status := <-ready:
		// TextToSpeech.SUCCESS is 0 and ERROR is -1.
		if status != 0 {
			s.Close()
			return nil, ErrNoEngine
		}
	case <-time.After(10 * time.Second):
		s.Close()
		return nil, errors.New("speech: the engine never said it was ready")
	}
	return s, nil
}
