//go:build darwin

package antui

/*
#cgo CFLAGS: -Wno-deprecated-declarations
#cgo LDFLAGS: -framework AudioToolbox -framework CoreFoundation
#include "antui_audio_darwin.h"
*/
import "C"

import "time"

// coreaudio is one open stream: the queue C owns and the mixer Go feeds it
// from.
type coreaudio struct {
	audio *Audio
	id    int
	queue C.AudioQueueRef
	delay time.Duration
}

func (c *coreaudio) latency() time.Duration { return c.delay }

func (c *coreaudio) close() error {
	queuesMu.Lock()
	delete(queues, c.id)
	queuesMu.Unlock()

	// Disposing waits for the callback to finish, so after this nothing is
	// looking the stream up any more.
	C.antuiAudioClose(c.queue)
	c.queue = nil
	return nil
}
