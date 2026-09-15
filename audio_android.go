//go:build android

package antui

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// Sound on Android goes through AAudio.
//
// The alternative below Android 8 is OpenSL ES — several hundred lines of
// COM-shaped C for a platform Google itself has deprecated it on, and for a
// share of devices now under one per cent. So an older device gets
// [ErrNoAudio], which is the same answer the rest of this library gives a
// machine with no sound card, and a game carries on in silence rather than
// refusing to start.
//
// The stream is **written to**, not pulled from. AAudio has a callback mode
// that is lower latency, and it runs the callback on a real-time thread —
// which Go's collector may stop, which is the one thing a real-time audio
// thread must never have happen to it. A goroutine that fills a block and
// blocks on the write is a few milliseconds slower and cannot glitch for
// that reason.
type androidAudio struct {
	stream *ndk.AudioStream

	stopOnce sync.Once
	stop     chan struct{}
	done     chan struct{}
}

func openAudio(name string, a *Audio) (driver, error) {
	s, err := ndk.OpenAudio(a.rate, a.channels, 0)
	if err != nil {
		if errors.Is(err, ndk.ErrNoAAudio) {
			return nil, fmt.Errorf("%w: %v", ErrNoAudio, err)
		}
		return nil, err
	}
	// What the device settled on, which is not always what was asked for:
	// a phone that runs its mixer at 48000 gives 48000 whatever the request.
	a.rate = s.Rate
	a.channels = s.Channels

	d := &androidAudio{
		stream: s,
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}
	go d.run(a)
	return d, nil
}

// blockBursts is how many of the device's own bursts go in one write.
//
// One is the smallest the device can take and leaves no room for a late
// goroutine; four is about ten milliseconds on a typical phone, which is
// under the threshold where a person hears the delay and far enough above
// the burst to survive a garbage collection.
const blockBursts = 4

func (d *androidAudio) run(a *Audio) {
	defer close(d.done)

	frames := d.stream.Burst * blockBursts
	if frames <= 0 {
		frames = 1024
	}
	block := make([]float32, frames*d.stream.Channels)

	for {
		select {
		case <-d.stop:
			return
		default:
		}
		a.pull(block)
		// A timeout rather than a wait for ever: a device that has been
		// taken away — a headset unplugged mid-write — otherwise leaves this
		// goroutine here for the life of the process.
		if _, err := d.stream.Write(block, time.Second); err != nil {
			if !a.stopped() {
				a.fail(err)
			}
			return
		}
	}
}

func (d *androidAudio) close() error {
	d.stopOnce.Do(func() {
		close(d.stop)
		// Wait for the goroutine to stop writing before the stream goes:
		// closing one that is being written to is a use-after-free inside
		// the platform.
		select {
		case <-d.done:
		case <-time.After(2 * time.Second):
		}
	})
	return d.stream.Close()
}

func (d *androidAudio) latency() time.Duration { return d.stream.Latency() }
