//go:build android

package sensor

import (
	"errors"
	"runtime"
	"time"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// run is the stream's own thread: it makes the looper, the queue and the
// subscriptions, and stays there until Close.
func (s *Stream) run(m *ndk.SensorManager, chosen []ndk.Sensor, opt Options, ready chan<- error) {
	// A looper belongs to its thread, and so does everything delivered to
	// it, so the goroutine may never move.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer close(s.stopped)

	looper := ndk.PrepareLooper()
	if looper == nil {
		ready <- errors.New("sensor: cannot make a looper for the sensor thread")
		return
	}
	queue, err := m.Queue(looper)
	if err != nil {
		ready <- err
		return
	}
	defer queue.Close()

	rate := int(opt.Rate.Microseconds())
	latency := opt.Latency.Microseconds()
	for _, sn := range chosen {
		// Batching first, and the plain path when the platform is too old
		// for it or the caller asked for none.
		if latency > 0 {
			if err := queue.Register(sn, rate, latency); err == nil {
				continue
			}
		}
		if err := queue.Enable(sn); err != nil {
			ready <- err
			return
		}
		if rate > 0 {
			if err := queue.SetRate(sn, rate); err != nil {
				ready <- err
				return
			}
		}
	}
	ready <- nil

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
