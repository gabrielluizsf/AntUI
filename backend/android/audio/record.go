//go:build android

package audio

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/backend/android/permission"
)

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
