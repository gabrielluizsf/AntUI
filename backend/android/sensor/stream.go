//go:build android

package sensor

import (
	"sync"
)

// Stream is a set of running sensors.
//
// It owns a thread of its own with a looper of its own, because that is what
// the platform's queue needs and because sensors run far faster than a frame
// loop — an accelerometer at its fastest is 200 readings a second, and
// pulling those through the frame loop would tie the two together for no
// reason.
type Stream struct {
	events chan Reading

	mu     sync.Mutex
	latest map[Kind]Reading
	fresh  map[Kind]bool

	stop     chan struct{}
	stopped  chan struct{}
	stopOnce sync.Once
}
