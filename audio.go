package antui

import (
	"errors"
	"sync"
	"time"
)

// Sound going out of the machine, on the same terms as the rest of AntUI:
// nothing to download, nothing to link against, and the system's own way of
// being spoken to rather than a library wrapping it.
//
//	Windows   →  WASAPI, the COM interfaces driven through syscall
//	macOS     →  CoreAudio's audio queues
//	Linux/BSD →  the PulseAudio protocol straight over the socket, which is
//	             what PipeWire answers on as well
//
// The device pulls: it calls fill whenever it needs more sound, from a
// goroutine of its own, and fill has to answer at once. Everything a game
// does about sound happens somewhere else and leaves the samples where the
// mixer can hand them over.

// ErrNoAudio is what opening gives back where there is no sound at all: a
// machine with no card, a session with no server, a platform nothing is
// written for yet. It is worth telling apart from a broken device, because a
// game carries on perfectly well in silence and should not stop for it.
var ErrNoAudio = errors.New("antui: no audio device")

// An Audio is one stream of sound going out to the speakers.
type Audio struct {
	rate     int
	channels int
	fill     func([]float32)

	mu   sync.Mutex
	err  error
	drv  driver
	done bool
}

// A driver is one platform's way of doing it.
type driver interface {
	close() error
	latency() time.Duration
}

// OpenAudio starts sending sound, and calls fill whenever the device wants
// more of it.
//
// fill is given interleaved stereo — left, right, left, right — at the rate
// asked for, and must overwrite all of it: what is in the slice is the last
// block, not silence. It is called from the device's own goroutine, so
// whatever it touches has to be safe to touch from there.
//
// The rate is what the game wants and not necessarily what it gets. Ask the
// device with Rate before making anything that depends on it.
func OpenAudio(name string, rate int, fill func([]float32)) (*Audio, error) {
	if fill == nil {
		return nil, errors.New("antui: audio with nothing to play")
	}
	if rate <= 0 {
		rate = 44100
	}
	if name == "" {
		name = "AntUI"
	}

	a := &Audio{rate: rate, channels: 2, fill: fill}
	drv, err := openAudio(name, a)
	if err != nil {
		return nil, err
	}
	a.drv = drv
	return a, nil
}

// Rate is how many frames a second the device is really running at.
func (a *Audio) Rate() int {
	if a == nil {
		return 0
	}
	return a.rate
}

// Channels is how many there are in a frame. Two, everywhere, for now.
func (a *Audio) Channels() int {
	if a == nil {
		return 0
	}
	return a.channels
}

// Latency is roughly how long it is between a sample being written and being
// heard. It is what the device was talked into rather than a measurement.
func (a *Audio) Latency() time.Duration {
	if a == nil || a.drv == nil {
		return 0
	}
	return a.drv.latency()
}

// Err is what went wrong after the device was opened — a card unplugged, a
// server that went away — or nil.
//
// A device that dies does not take the game with it. Everything goes on being
// mixed and nothing is heard, which is what happens on a machine with the
// speakers turned off anyway.
func (a *Audio) Err() error {
	if a == nil {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.err
}

// fail records why the device stopped. The first reason is kept: what follows
// a broken pipe is a hundred more broken pipes.
func (a *Audio) fail(err error) {
	if a == nil || err == nil {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.err == nil {
		a.err = err
	}
}

// stopped reports whether Close has been called, which is what tells a
// backend's goroutine that the socket closing under it was on purpose.
func (a *Audio) stopped() bool {
	if a == nil {
		return true
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.done
}

// Close stops the sound. It is safe to call twice.
func (a *Audio) Close() error {
	if a == nil {
		return nil
	}
	a.mu.Lock()
	if a.done {
		a.mu.Unlock()
		return nil
	}
	a.done = true
	drv := a.drv
	a.mu.Unlock()

	if drv == nil {
		return nil
	}
	return drv.close()
}

// pull fills a block of frames, and silences it when there is nothing to
// pull from. A backend hands its own buffer in.
func (a *Audio) pull(block []float32) {
	if a == nil || a.fill == nil {
		for i := range block {
			block[i] = 0
		}
		return
	}
	a.fill(block)
}