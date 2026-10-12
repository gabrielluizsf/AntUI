//go:build android

package speech

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// stop destroys the recogniser and releases the listeners. It is safe to
// call more than once.
func (l *listenState) stop(release func()) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			l.mu.Lock()
			l.stopped = true
			l.mu.Unlock()
			app.RunOnUISync(func() {
				jni.Do(func(e *jni.Env) error {
					if !l.recogniser.IsNil() {
						// destroy rather than stopListening: the recogniser
						// holds the microphone until it is destroyed, and a
						// microphone left open is one the user's indicator
						// keeps showing.
						e.InvokeVoid(l.recogniser, "destroy", jni.Sig(jni.TVoid))
						e.DeleteGlobal(l.recogniser)
					}
					if !l.intent.IsNil() {
						e.DeleteGlobal(l.intent)
					}
					return nil
				})
			})
			release()
		})
	}
}
