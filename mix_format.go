//go:build windows

package antui

import (
	"fmt"
	"unsafe"
)

// mixFormat returns the shared mixer's format. IAudioClient::GetMixFormat is
// the ninth method. What comes back is what the shared mixer is running at,
// and taking it as it stands is what makes this work on a device at any rate
// at all — the mixer above follows, rather than a resampler being written
// here.
func (w *wasapi) mixFormat(client uintptr) (*waveFormatExtensible, error) {
	var mix *waveFormatExtensible
	if err := failed(comCall(client, 8, uintptr(unsafe.Pointer(&mix))),
		"GetMixFormat"); err != nil {
		return nil, err
	}
	defer procCoTaskMemFree.Call(uintptr(unsafe.Pointer(mix)))

	// A copy, because the original is COM's memory and is about to go.
	format := new(waveFormatExtensible)
	*format = *mix
	if !w.understands(format) {
		return nil, fmt.Errorf(
			"antui: the device wants %d-bit samples in format %d, which this does not write",
			format.bitsPerSample, format.tag)
	}
	return format, nil
}
