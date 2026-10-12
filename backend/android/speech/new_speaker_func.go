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
	initToken, releaseInit := app.Listen(initListener(ready))
	s.release = append(s.release, releaseInit)

	progressToken, releaseProgress := app.Listen(s.onProgress)
	s.release = append(s.release, releaseProgress)

	var initErr error
	if e := app.RunOnUISync(func() {
		initErr = jni.Do(func(e *jni.Env) error {
			return makeEngine(e, s, engine, initToken, progressToken)
		})
	}); e != nil {
		s.Close()
		return nil, e
	}
	if initErr != nil {
		s.Close()
		return nil, initErr
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
