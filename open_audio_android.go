//go:build android

package antui

import (
	"errors"
	"fmt"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

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
