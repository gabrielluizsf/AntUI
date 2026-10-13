//go:build windows

package antui

import "fmt"

// feed is the stream's loop: wait to be asked, hand over a block, repeat.
func (w *wasapi) feed(client, render uintptr, frames uint32, format *waveFormatExtensible) {
	s := wasapiStream{
		client:   client,
		render:   render,
		frames:   frames,
		channels: int(format.channels),
		asFloats: w.floats(format),
		block:    make([]float32, int(frames)*2),
	}

	for {
		state, _, _ := procWaitForSingleObject.Call(w.event, 2000)
		select {
		case <-w.quit:
			return
		default:
		}
		if state == waitTimeout {
			// Two seconds without the device asking for anything is a device
			// that has gone away — pulled out, or switched to another one.
			w.audio.fail(fmt.Errorf("antui: the audio device stopped asking for sound"))
			return
		}
		if state != waitObject0 {
			w.audio.fail(fmt.Errorf("antui: waiting on the audio device gave %d", state))
			return
		}
		if !w.feedOnce(s) {
			return
		}
	}
}
