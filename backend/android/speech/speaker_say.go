//go:build android

package speech

import (
	"errors"
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

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
