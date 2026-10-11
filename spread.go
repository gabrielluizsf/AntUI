//go:build windows

package antui

import (
	"unsafe"
)

// spread writes a stereo block into the device's buffer, in the device's own
// layout.
//
// The mixer above is stereo and a device may be anything: a headset is two, a
// living room is six. The two go to the front pair and the rest are left
// silent, which is what every game does with a stereo mix — putting the same
// sound in the surrounds as well makes a room sound like a corridor.
func (w *wasapi) spread(buffer uintptr, block []float32, frames, channels int, asFloats bool) {
	if asFloats {
		out := unsafe.Slice((*float32)(unsafe.Pointer(buffer)), frames*channels)
		for i := range out {
			out[i] = 0
		}
		for frame := range frames {
			out[frame*channels] = block[frame*2]
			if channels > 1 {
				out[frame*channels+1] = block[frame*2+1]
			}
		}
		return
	}

	out := unsafe.Slice((*int16)(unsafe.Pointer(buffer)), frames*channels)
	for i := range out {
		out[i] = 0
	}
	for frame := range frames {
		out[frame*channels] = pcm16(block[frame*2])
		if channels > 1 {
			out[frame*channels+1] = pcm16(block[frame*2+1])
		}
	}
}
