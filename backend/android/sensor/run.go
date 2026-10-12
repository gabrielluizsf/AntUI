//go:build android

package sensor

import (
	"errors"
	"runtime"

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
	if err := registerAll(queue, chosen, rate, latency); err != nil {
		ready <- err
		return
	}
	ready <- nil
	pollLoop(s, looper, queue)
}
