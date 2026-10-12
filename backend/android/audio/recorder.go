//go:build android

package audio

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// Recorder is sound coming in from the microphone.
type Recorder struct {
	stream *ndk.AudioStream

	stopOnce sync.Once
	stop     chan struct{}
	done     chan struct{}

	mu  sync.Mutex
	err error
}
