//go:build android

package audio

import (
	"errors"
	"sync"
	"time"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/backend/android/permission"
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

// ErrNoMicrophone is recording without the permission having been granted.
var ErrNoMicrophone = errors.New("audio: this app has not been granted " +
	"RECORD_AUDIO; ask for it with antui/backend/android/permission first")

// Record starts listening, and calls f with each block that arrives.
//
// f runs on a goroutine of the recorder's own and must answer at once —
// whatever it does with the samples, it should be copying them somewhere and
// nothing more. The slice it is given is reused.
//
// **A recording app has to say so.** The permission is only half of it: from
// Android 9 the microphone is cut off entirely while the app is not in
// front, and from Android 12 an indicator appears whenever it is live. An
// app that records without making that plain is one the store removes.
func Record(rate, channels int, f func([]float32)) (*Recorder, error) {
	if f == nil {
		return nil, errors.New("audio: recording with nowhere to put it")
	}
	held, err := permission.Held(permission.Microphone)
	if err != nil {
		return nil, err
	}
	if !held {
		return nil, ErrNoMicrophone
	}
	if rate <= 0 {
		rate = 44100
	}
	if channels <= 0 {
		channels = 1 // a microphone is one, whatever is asked for
	}

	s, err := ndk.OpenAudioInput(rate, channels, 0)
	if err != nil {
		return nil, err
	}
	r := &Recorder{
		stream: s,
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}
	go r.run(f)
	return r, nil
}

func (r *Recorder) run(f func([]float32)) {
	defer close(r.done)

	frames := r.stream.Burst * blockBursts
	if frames <= 0 {
		frames = 1024
	}
	block := make([]float32, frames*r.stream.Channels)

	for {
		select {
		case <-r.stop:
			return
		default:
		}
		n, err := r.stream.Read(block, time.Second)
		if err != nil {
			r.mu.Lock()
			if r.err == nil {
				r.err = err
			}
			r.mu.Unlock()
			return
		}
		if n > 0 {
			f(block[:n*r.stream.Channels])
		}
	}
}

// blockBursts is how many of the device's own bursts go in one read, for the
// same reason it is what it is on the way out.
const blockBursts = 4

// Rate is how many frames a second the device is really recording at.
func (r *Recorder) Rate() int { return r.stream.Rate }

// Channels is how many the device is really giving.
func (r *Recorder) Channels() int { return r.stream.Channels }

// Err is what stopped it, or nil.
func (r *Recorder) Err() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.err
}

// Close stops recording.
func (r *Recorder) Close() error {
	r.stopOnce.Do(func() {
		close(r.stop)
		select {
		case <-r.done:
		case <-time.After(2 * time.Second):
		}
	})
	return r.stream.Close()
}
