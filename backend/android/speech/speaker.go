//go:build android

package speech

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

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
