//go:build windows

package antui

import (
	"time"
)

// start opens the default endpoint and gets as far as a running stream. It
// answers the client, the render interface, how many frames the device's
// buffer holds, and the format it wants them in.
func (w *wasapi) start() (client, render uintptr, frames uint32, format *waveFormatExtensible, err error) {
	device, err := w.defaultDevice()
	if err != nil {
		return 0, 0, 0, nil, err
	}
	defer comRelease(device)

	client, err = w.activate(device)
	if err != nil {
		return 0, 0, 0, nil, err
	}
	format, err = w.mixFormat(client)
	if err != nil {
		comRelease(client)
		return 0, 0, 0, nil, err
	}
	render, frames, err = w.initialize(client, format)
	if err != nil {
		comRelease(client)
		return 0, 0, 0, nil, err
	}

	w.audio.rate = int(format.samplesPerSec)
	w.delay = time.Duration(frames) * time.Second / time.Duration(format.samplesPerSec)
	return client, render, frames, format, nil
}
