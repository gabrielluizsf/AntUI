//go:build android

package speech

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

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
