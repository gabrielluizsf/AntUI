//go:build android

package sensor

import (
	"time"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// pollLoop waits on the looper and delivers everything the queue has.
func pollLoop(s *Stream, looper *ndk.Looper, queue *ndk.SensorQueue) {
	buf := make([]ndk.SensorEvent, 32)
	for {
		select {
		case <-s.stop:
			return
		default:
		}
		// A timeout rather than a wait, so that a Close that races with the
		// poll is noticed even if the wake is missed.
		looper.Poll(250)
		for {
			n := queue.Read(buf)
			if n == 0 {
				break
			}
			for _, e := range buf[:n] {
				s.deliver(Reading{
					Kind: Kind(e.Type),
					At:   time.Duration(e.At),
					V:    e.V,
				})
			}
		}
	}
}
