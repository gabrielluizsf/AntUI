package antui

import (
	"errors"
)

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
