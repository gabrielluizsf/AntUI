//go:build android

package speech

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

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
