//go:build android

package speech

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// listenState is what Listen holds across the three threads it lives on: the
// caller's, the UI thread the engine answers on, and the stop that may come
// from either.
type listenState struct {
	opt        Options
	report     func(Heard, error)
	recogniser jni.Object
	intent     jni.Object
	stopped    bool
	mu         sync.Mutex
}
