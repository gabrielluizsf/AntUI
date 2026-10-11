//go:build darwin

package antui

/*
#cgo CFLAGS: -Wno-deprecated-declarations
#cgo LDFLAGS: -framework AudioToolbox -framework CoreFoundation
#include "antui_audio_darwin.h"
*/
import "C"

import (
	"fmt"
	"time"
)

// openAudio is the macOS half of OpenAudio.
//
// CoreAudio through an audio queue, which is the part of it that takes
// samples in the format they are already in and deals with the hardware's own
// ideas about rate and layout on the way past.
//
// The lower AudioUnit is what a professional tool reaches for; the queue is a
// tenth of the code and its extra few milliseconds are under what a screen
// refresh costs anyway.
//
// NOT YET COMPILED ON macOS. Like the Cocoa half of the window backend, this
// is written from the framework headers and waits for a machine to run it —
// see audio_darwin_test.go, which only runs there.
func openAudio(name string, a *Audio) (driver, error) {
	// Ten milliseconds a buffer and three of them: a third of a frame at
	// sixty, and short enough that a sound lands with the thing that made it.
	const buffers = 3
	frames := max(a.rate/100, 64)

	c := &coreaudio{audio: a}

	queuesMu.Lock()
	nextQueue++
	c.id = nextQueue
	queues[c.id] = c
	queuesMu.Unlock()

	var queue C.AudioQueueRef
	status := C.antuiAudioOpen(C.double(a.rate), C.int(c.id), C.int(frames),
		C.int(buffers), &queue)
	if status != 0 {
		queuesMu.Lock()
		delete(queues, c.id)
		queuesMu.Unlock()
		return nil, fmt.Errorf("%w: CoreAudio refused a stream (%d)", ErrNoAudio, int(status))
	}

	c.queue = queue
	c.delay = time.Duration(frames*buffers) * time.Second / time.Duration(a.rate)
	return c, nil
}
