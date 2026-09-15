//go:build android

package speech

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// How a new piece of speech meets one already being said.
const (
	// Replace stops whatever is being said and says this instead. It is what
	// almost everything wants: the thing being said is out of date the
	// moment there is something newer to say.
	Replace int32 = 0
	// Queue waits its turn.
	Queue int32 = 1
)

// ErrNoEngine is a device with no speech engine installed. It is not rare:
// the engine is an app, and a device without Google's services may have none.
var ErrNoEngine = errors.New("speech: this device has no text-to-speech engine")

// Speaker turns text into sound.
//
// Making one takes a moment — the engine is another app and has to be bound
// to — so it is worth making one and keeping it rather than one per
// sentence. Closing it matters: an engine left bound is a service the system
// keeps alive for an app that has stopped caring.
type Speaker struct {
	mu      sync.Mutex
	tts     jni.Object
	closed  bool
	release []func()

	// done is closed when a particular utterance finishes, so that a caller
	// can wait for one.
	waitMu sync.Mutex
	waits  map[string]chan error
	next   int
}

// NewSpeaker starts the user's own engine and waits for it to be ready.
//
// It waits because there is nothing useful to do with an engine that is not:
// asking it to speak before it has bound is refused, silently, and the words
// are simply never said.
//
// A device with an engine installed but none chosen as the default has no
// default, and this fails with [ErrNoEngine]. [Engines] lists what is there
// and [NewSpeakerWith] names one.
func NewSpeaker() (*Speaker, error) { return newSpeaker("") }

// NewSpeakerWith starts a particular engine, named by its package.
//
// Prefer [NewSpeaker]: which engine speaks is the user's choice, made in the
// settings, and an app that overrides it overrides an accessibility setting.
// This is for the case where there is no choice to respect — no default set
// — and for a test that has to be sure which engine answered.
func NewSpeakerWith(engine string) (*Speaker, error) { return newSpeaker(engine) }

// Engines is every text-to-speech engine installed, by package name. The
// list is empty on a device with none.
func Engines() ([]string, error) {
	var out []string
	err := jni.Do(func(e *jni.Env) error {
		// A throwaway engine, only to ask it what else is there. It is the
		// only way the platform offers: the list is an instance method.
		tts, err := e.Make("android/speech/tts/TextToSpeech",
			jni.Sig(jni.TVoid, jni.TContext,
				jni.TClass("android/speech/tts/TextToSpeech$OnInitListener")),
			jni.Ref(app.Context()), jni.Ref(jni.Object{}))
		if err != nil {
			return err
		}
		defer e.InvokeVoid(tts, "shutdown", jni.Sig(jni.TVoid))

		list, err := e.Invoke(tts, "getEngines", jni.Sig(jni.TClass("java/util/List")))
		if err != nil || list.IsNil() {
			return err
		}
		n, err := e.InvokeInt(list, "size", jni.Sig(jni.TInt))
		if err != nil {
			return err
		}
		return e.Frame(int(n)*4+16, func() error {
			for i := range int(n) {
				info, err := e.Invoke(list, "get", jni.Sig(jni.TObject, jni.TInt),
					jni.Int(int32(i)))
				if err != nil || info.IsNil() {
					continue
				}
				cls, err := e.Class("android/speech/tts/TextToSpeech$EngineInfo")
				if err != nil {
					return err
				}
				f, err := e.Field(cls, "name", jni.TString)
				if err != nil {
					return err
				}
				name, err := e.GetString(info, f)
				if err == nil && name != "" {
					out = append(out, name)
				}
			}
			return nil
		})
	})
	return out, err
}

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

// onProgress is what the shim's listener reports.
func (s *Speaker) onProgress(e *jni.Env, method string, args []jni.Object) {
	if len(args) == 0 {
		return
	}
	id, err := e.GoString(args[0])
	if err != nil || id == "" {
		return
	}
	var answer error
	switch method {
	case "onDone":
	case "onError":
		answer = errors.New("speech: the engine could not say it")
	default:
		return // onStart, which nothing is waiting for
	}
	s.waitMu.Lock()
	ch := s.waits[id]
	delete(s.waits, id)
	s.waitMu.Unlock()
	if ch != nil {
		ch <- answer
	}
}

// Say speaks, and returns as soon as the words are queued.
func (s *Speaker) Say(text string, mode int32) error {
	_, err := s.say(text, mode, false)
	return err
}

// SayAndWait speaks and waits until the words have been said.
//
// **Not from the frame loop**: it takes as long as the sentence does.
func (s *Speaker) SayAndWait(text string, mode int32) error {
	done, err := s.say(text, mode, true)
	if err != nil {
		return err
	}
	return <-done
}

func (s *Speaker) say(text string, mode int32, wait bool) (chan error, error) {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil, errors.New("speech: the speaker is closed")
	}
	s.next++
	id := fmt.Sprintf("antui-%d", s.next)
	tts := s.tts
	s.mu.Unlock()

	var done chan error
	if wait {
		done = make(chan error, 1)
		s.waitMu.Lock()
		s.waits[id] = done
		s.waitMu.Unlock()
	}

	err := jni.Do(func(e *jni.Env) error {
		jtext, err := e.String(text)
		if err != nil {
			return err
		}
		jid, err := e.String(id)
		if err != nil {
			return err
		}
		r, err := e.InvokeInt(tts, "speak",
			jni.Sig(jni.TInt, jni.TClass("java/lang/CharSequence"), jni.TInt,
				jni.TClass("android/os/Bundle"), jni.TString),
			jni.Ref(jtext), jni.Int(mode), jni.Ref(jni.Object{}), jni.Ref(jid))
		if err != nil {
			return err
		}
		if r != 0 {
			return errors.New("speech: the engine refused the words")
		}
		return nil
	})
	if err != nil {
		s.waitMu.Lock()
		delete(s.waits, id)
		s.waitMu.Unlock()
		return nil, err
	}
	return done, nil
}

// Speaking reports whether anything is being said now.
func (s *Speaker) Speaking() (bool, error) {
	var out bool
	err := jni.Do(func(e *jni.Env) error {
		var err error
		out, err = e.InvokeBool(s.tts, "isSpeaking", jni.Sig(jni.TBool))
		return err
	})
	return out, err
}

// Stop drops everything queued and stops mid-word.
func (s *Speaker) Stop() error {
	return jni.Do(func(e *jni.Env) error {
		_, err := e.InvokeInt(s.tts, "stop", jni.Sig(jni.TInt))
		return err
	})
}

// SetRate is how fast: 1 is normal, 2 is twice as fast, 0.5 is half.
func (s *Speaker) SetRate(rate float32) error { return s.tune("setSpeechRate", rate) }

// SetPitch is how high: 1 is normal.
func (s *Speaker) SetPitch(pitch float32) error { return s.tune("setPitch", pitch) }

func (s *Speaker) tune(method string, v float32) error {
	return jni.Do(func(e *jni.Env) error {
		_, err := e.InvokeInt(s.tts, method, jni.Sig(jni.TInt, jni.TFloat), jni.Float(v))
		return err
	})
}

// The answers setLanguage gives.
const (
	langMissingData  = -1
	langNotSupported = -2
)

// SetLanguage asks for a language by its tag — "pt-BR", "en-GB".
//
// It reports whether the engine has it. A language it knows but has not
// downloaded is as unusable as one it does not know, and both come back
// false: the difference is not something an app can do anything about.
func (s *Speaker) SetLanguage(tag string) (bool, error) {
	var ok bool
	err := jni.Do(func(e *jni.Env) error {
		jtag, err := e.String(tag)
		if err != nil {
			return err
		}
		locale, err := e.Static("java/util/Locale", "forLanguageTag",
			jni.Sig(jni.TClass("java/util/Locale"), jni.TString), jni.Ref(jtag))
		if err != nil {
			return err
		}
		r, err := e.InvokeInt(s.tts, "setLanguage",
			jni.Sig(jni.TInt, jni.TClass("java/util/Locale")), jni.Ref(locale))
		if err != nil {
			return err
		}
		ok = r != langMissingData && r != langNotSupported
		return nil
	})
	return ok, err
}

// Close stops the engine and lets go of it.
func (s *Speaker) Close() error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true
	tts := s.tts
	releases := s.release
	s.tts = jni.Object{}
	s.mu.Unlock()

	if !tts.IsNil() {
		jni.Do(func(e *jni.Env) error {
			e.InvokeVoid(tts, "shutdown", jni.Sig(jni.TVoid))
			e.DeleteGlobal(tts)
			return nil
		})
	}
	for _, r := range releases {
		r()
	}
	// Anything still waiting is never going to be answered.
	s.waitMu.Lock()
	for id, ch := range s.waits {
		ch <- errors.New("speech: the speaker was closed")
		delete(s.waits, id)
	}
	s.waitMu.Unlock()
	return nil
}
