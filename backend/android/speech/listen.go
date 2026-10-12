//go:build android

package speech

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/permission"
)

// Listen starts the recogniser and reports what it hears.
//
// f runs on the UI thread, which is where the engine answers. Every call has
// either a result or an error, never both.
//
// It needs RECORD_AUDIO. On almost every device the sound goes to a server to
// be recognised, so an app that listens is an app that sends what was said
// somewhere — which the user should be told, in the app and not only in a
// privacy policy.
func Listen(opt Options, f func(Heard, error)) (stop func(), err error) {
	held, err := permission.Held(permission.Microphone)
	if err != nil {
		return nil, err
	}
	if !held {
		return nil, ErrNoPermission
	}
	if !Available() {
		return nil, ErrNoRecogniser
	}

	var (
		recogniser jni.Object
		intent     jni.Object
		stopped    bool
		mu         sync.Mutex
	)

	token, release := app.Listen(func(e *jni.Env, method string, args []jni.Object) {
		switch method {
		case "onResults", "onPartialResults":
			if len(args) == 0 {
				return
			}
			heard, ok := readResults(e, args[0], method == "onPartialResults")
			if ok {
				f(heard, nil)
			}
			mu.Lock()
			again := opt.Continuous && !stopped && method == "onResults"
			mu.Unlock()
			if again {
				// Already on the UI thread, which is where startListening
				// has to be called from.
				startListening(e, recogniser, intent)
			}
		case "onError":
			if len(args) == 0 {
				return
			}
			code, err := e.InvokeInt(args[0], "intValue", jni.Sig(jni.TInt))
			if err != nil {
				return
			}
			f(Heard{}, recogniserError(code))
		}
	})

	fail := func(err error) (func(), error) {
		release()
		return nil, err
	}
	var setupErr error
	if e := app.RunOnUISync(func() {
		setupErr = jni.Do(func(e *jni.Env) error {
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
			recogniser = e.Global(r)

			listener, err := app.Proxy(e, token, "android.speech.RecognitionListener")
			if err != nil {
				return err
			}
			if err := e.InvokeVoid(recogniser, "setRecognitionListener",
				jni.Sig(jni.TVoid, jni.TClass("android/speech/RecognitionListener")),
				jni.Ref(listener)); err != nil {
				return err
			}

			in, err := buildIntent(e, opt)
			if err != nil {
				return err
			}
			intent = e.Global(in)
			return startListening(e, recogniser, intent)
		})
	}); e != nil {
		return fail(e)
	}
	if setupErr != nil {
		return fail(setupErr)
	}

	var once sync.Once
	return func() {
		once.Do(func() {
			mu.Lock()
			stopped = true
			mu.Unlock()
			app.RunOnUISync(func() {
				jni.Do(func(e *jni.Env) error {
					if !recogniser.IsNil() {
						// destroy rather than stopListening: the recogniser
						// holds the microphone until it is destroyed, and a
						// microphone left open is one the user's indicator
						// keeps showing.
						e.InvokeVoid(recogniser, "destroy", jni.Sig(jni.TVoid))
						e.DeleteGlobal(recogniser)
					}
					if !intent.IsNil() {
						e.DeleteGlobal(intent)
					}
					return nil
				})
			})
			release()
		})
	}, nil
}
