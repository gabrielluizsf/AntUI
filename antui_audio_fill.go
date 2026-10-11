//go:build darwin

package antui

/*
#cgo CFLAGS: -Wno-deprecated-declarations
#cgo LDFLAGS: -framework AudioToolbox -framework CoreFoundation
#include "antui_audio_darwin.h"
*/
import "C"

import "unsafe"

//export antuiAudioFill
func antuiAudioFill(id C.int, samples *C.float, frames C.int) {
	// The buffer belongs to CoreAudio, so it is written through rather than
	// copied into: a slice over C memory, which is exactly the shape the
	// mixer writes.
	fillQueue(int(id), unsafe.Slice((*float32)(unsafe.Pointer(samples)), int(frames)*2))
}
