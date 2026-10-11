//go:build windows

package antui

import (
	"fmt"
	"unsafe"
)

// initialize sets the stream up and starts it.
func (w *wasapi) initialize(client uintptr, format *waveFormatExtensible) (render uintptr, frames uint32, err error) {
	// IAudioClient::Initialize, the fourth. A duration of zero in shared mode
	// asks for the engine's own period, which is the shortest it will give.
	if err := failed(comCall(client, 3, shareModeShared, streamEventCallback,
		0, 0, uintptr(unsafe.Pointer(&format.waveFormatEx)), 0),
		"Initialize"); err != nil {
		return 0, 0, err
	}

	event, _, _ := procCreateEventW.Call(0, 0, 0, 0)
	if event == 0 {
		return 0, 0, fmt.Errorf("antui: no event for the audio device")
	}
	w.event = event

	// SetEventHandle is the fourteenth.
	if err := failed(comCall(client, 13, event), "SetEventHandle"); err != nil {
		return 0, 0, err
	}
	// GetBufferSize, the fifth.
	if err := failed(comCall(client, 4, uintptr(unsafe.Pointer(&frames))),
		"GetBufferSize"); err != nil {
		return 0, 0, err
	}
	// GetService, the fifteenth.
	if err := failed(comCall(client, 14,
		uintptr(unsafe.Pointer(&iidIAudioRenderClient)),
		uintptr(unsafe.Pointer(&render))), "GetService"); err != nil {
		return 0, 0, err
	}
	// Start, the eleventh.
	if err := failed(comCall(client, 10), "Start"); err != nil {
		comRelease(render)
		return 0, 0, err
	}
	return render, frames, nil
}
