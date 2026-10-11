//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"sync"
)

// Sound on Android goes through AAudio.
//
// The alternative below Android 8 is OpenSL ES — several hundred lines of
// COM-shaped C for a platform Google itself has deprecated it on, and for a
// share of devices now under one per cent. So an older device gets
// [ErrNoAudio], which is the same answer the rest of this library gives a
// machine with no sound card, and a game carries on in silence rather than
// refusing to start.
//
// The stream is **written to**, not pulled from. AAudio has a callback mode
// that is lower latency, and it runs the callback on a real-time thread —
// which Go's collector may stop, which is the one thing a real-time audio
// thread must never have happen to it. A goroutine that fills a block and
// blocks on the write is a few milliseconds slower and cannot glitch for
// that reason.
type androidAudio struct {
	stream *ndk.AudioStream

	stopOnce sync.Once
	stop     chan struct{}
	done     chan struct{}
}
