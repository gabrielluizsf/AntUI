//go:build windows

package antui

import "unsafe"

// feedOnce fills whatever room the device has and hands it back.
func (w *wasapi) feedOnce(s wasapiStream) bool {
	// GetCurrentPadding, the seventh: how much of the buffer the device has
	// not played yet. The rest is ours to fill.
	var padding uint32
	if err := failed(comCall(s.client, 6, uintptr(unsafe.Pointer(&padding))),
		"GetCurrentPadding"); err != nil {
		w.audio.fail(err)
		return false
	}
	want := s.frames - padding
	if want == 0 {
		return true
	}

	// IAudioRenderClient::GetBuffer, the fourth.
	var buffer uintptr
	if err := failed(comCall(s.render, 3, uintptr(want),
		uintptr(unsafe.Pointer(&buffer))), "GetBuffer"); err != nil {
		w.audio.fail(err)
		return false
	}

	w.audio.pull(s.block[:int(want)*2])
	w.spread(buffer, s.block[:int(want)*2], int(want), s.channels, s.asFloats)

	// ReleaseBuffer, the fifth.
	if err := failed(comCall(s.render, 4, uintptr(want), 0), "ReleaseBuffer"); err != nil {
		w.audio.fail(err)
		return false
	}
	return true
}
