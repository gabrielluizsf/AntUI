//go:build android

package speech

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/permission"
)

// Heard is what the recogniser made of what was said.
type Heard struct {
	// Text is its best guess.
	Text string
	// Alternatives are all of them, best first, with Text at the front.
	Alternatives []string
	// Confidence is 0 to 1 for Text, and 0 when the engine will not say —
	// which many will not, so it cannot be used as a threshold.
	Confidence float32
	// Partial says this is a guess so far and not the answer. More will
	// follow, and the last one will not be partial.
	Partial bool
}

// The reasons recognition stops, translated out of the engine's numbers.
var (
	ErrNoMatch      = errors.New("speech: nothing recognisable was said")
	ErrTimeout      = errors.New("speech: nothing was said")
	ErrNoNetwork    = errors.New("speech: recognition needs a network and there is none")
	ErrBusy         = errors.New("speech: the recogniser is busy")
	ErrNoPermission = errors.New("speech: this app has not been granted RECORD_AUDIO")
)

// ErrNoRecogniser is a device with nothing installed that can listen. It is
// common: recognition is a service, usually Google's, and a device without
// those services has none.
var ErrNoRecogniser = errors.New("speech: this device has nothing that can listen")

// Available reports whether anything on this device can listen.
func Available() bool {
	var out bool
	jni.Do(func(e *jni.Env) error {
		c, err := e.Class("android/speech/SpeechRecognizer")
		if err != nil {
			return err
		}
		m, err := e.StaticMethod(c, "isRecognitionAvailable",
			jni.Sig(jni.TBool, jni.TContext))
		if err != nil {
			return err
		}
		out, err = e.CallStaticBool(c, m, jni.Ref(app.Context()))
		return err
	})
	return out
}

// Options is how to listen.
type Options struct {
	// Language is a tag — "pt-BR", "en-GB". Empty uses the device's own.
	Language string
	// Partial reports guesses as they are made, not only the final answer.
	// They arrive with [Heard.Partial] set and are worth showing: a person
	// watching words appear knows they are being heard.
	Partial bool
	// Continuous starts listening again as soon as an answer is given, so
	// that a conversation can carry on. Without it, listening stops after
	// one utterance.
	//
	// It is not free — the microphone stays open — and on most devices the
	// engine is a network service, so it is also not private. Say so.
	Continuous bool
}

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

func startListening(e *jni.Env, recogniser, intent jni.Object) error {
	return e.InvokeVoid(recogniser, "startListening",
		jni.Sig(jni.TVoid, jni.TClass("android/content/Intent")), jni.Ref(intent))
}

func buildIntent(e *jni.Env, opt Options) (jni.Object, error) {
	action, err := e.ConstantString("android/speech/RecognizerIntent", "ACTION_RECOGNIZE_SPEECH")
	if err != nil {
		return jni.Object{}, err
	}
	jaction, err := e.String(action)
	if err != nil {
		return jni.Object{}, err
	}
	in, err := e.Make("android/content/Intent", jni.Sig(jni.TVoid, jni.TString),
		jni.Ref(jaction))
	if err != nil {
		return jni.Object{}, err
	}
	put := func(key, value string) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		v, err := e.String(value)
		if err != nil {
			return err
		}
		_, err = e.Invoke(in, "putExtra",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TString),
			jni.Ref(k), jni.Ref(v))
		return err
	}
	// The free-form model rather than the web-search one: the second is
	// tuned for short queries and gets a sentence wrong.
	if err := put("android.speech.extra.LANGUAGE_MODEL", "free_form"); err != nil {
		return jni.Object{}, err
	}
	if opt.Language != "" {
		if err := put("android.speech.extra.LANGUAGE", opt.Language); err != nil {
			return jni.Object{}, err
		}
	}
	if opt.Partial {
		k, err := e.String("android.speech.extra.PARTIAL_RESULTS")
		if err != nil {
			return jni.Object{}, err
		}
		if _, err := e.Invoke(in, "putExtra",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TBool),
			jni.Ref(k), jni.Bool(true)); err != nil {
			return jni.Object{}, err
		}
	}
	return in, nil
}

// readResults pulls the guesses out of the bundle the engine answers with.
func readResults(e *jni.Env, bundle jni.Object, partial bool) (Heard, bool) {
	if bundle.IsNil() {
		return Heard{}, false
	}
	key, err := e.String("results_recognition")
	if err != nil {
		return Heard{}, false
	}
	list, err := e.Invoke(bundle, "getStringArrayList",
		jni.Sig(jni.TClass("java/util/ArrayList"), jni.TString), jni.Ref(key))
	if err != nil || list.IsNil() {
		return Heard{}, false
	}
	n, err := e.InvokeInt(list, "size", jni.Sig(jni.TInt))
	if err != nil || n == 0 {
		return Heard{}, false
	}
	out := Heard{Partial: partial, Alternatives: make([]string, 0, n)}
	err = e.Frame(int(n)*4+16, func() error {
		for i := range int(n) {
			item, err := e.Invoke(list, "get", jni.Sig(jni.TObject, jni.TInt),
				jni.Int(int32(i)))
			if err != nil || item.IsNil() {
				continue
			}
			s, err := e.GoString(item)
			if err == nil && s != "" {
				out.Alternatives = append(out.Alternatives, s)
			}
		}
		return nil
	})
	if err != nil || len(out.Alternatives) == 0 {
		return Heard{}, false
	}
	out.Text = out.Alternatives[0]

	// The confidence scores, which many engines simply do not send.
	if ckey, err := e.String("confidence_scores"); err == nil {
		if scores, err := e.Invoke(bundle, "getFloatArray",
			jni.Sig(jni.TArray(jni.TFloat), jni.TString), jni.Ref(ckey)); err == nil &&
			!scores.IsNil() {
			if got, err := goFloats(e, scores); err == nil && len(got) > 0 {
				out.Confidence = got[0]
			}
		}
	}
	return out, true
}

// recogniserError turns the engine's number into something a person could
// act on.
func recogniserError(code int32) error {
	switch code {
	case 1, 2:
		return ErrNoNetwork
	case 6:
		return ErrTimeout
	case 7:
		return ErrNoMatch
	case 8:
		return ErrBusy
	case 9:
		return ErrNoPermission
	}
	return fmt.Errorf("speech: the recogniser stopped (%d)", code)
}
