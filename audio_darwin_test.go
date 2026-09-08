//go:build darwin && cgo

package antui

import (
	"errors"
	"testing"
	"time"
)

// These run on macOS and nowhere else: the audio queue, and the one piece of
// it that C owns — the number a stream is looked up by when the queue calls
// back into Go.

func TestCoreAudioPlaysOnThisMachine(t *testing.T) {
	var frames int
	done := make(chan struct{})
	closed := false

	audio, err := OpenAudio("AntUI test", 44100, func(block []float32) {
		for i := range block {
			block[i] = 0
		}
		frames += len(block) / 2
		if frames >= 4410 && !closed {
			closed = true
			close(done)
		}
	})
	if errors.Is(err, ErrNoAudio) {
		t.Skip("no sound on this machine:", err)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer audio.Close()

	t.Logf("the queue is at %d Hz with %v of latency", audio.Rate(), audio.Latency())
	if got := audio.Latency(); got <= 0 || got > 100*time.Millisecond {
		t.Errorf("the latency is %v", got)
	}

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatalf("the queue asked for %d frames in three seconds", frames)
	}
	if err := audio.Err(); err != nil {
		t.Fatal(err)
	}
}

// Closing takes the stream out of the table C looks it up in, and a queue
// that has already been asked for one more buffer must find nothing there
// rather than a stream that is halfway gone.
func TestCoreAudioForgetsAClosedStream(t *testing.T) {
	audio, err := OpenAudio("AntUI test", 44100, func(block []float32) {})
	if errors.Is(err, ErrNoAudio) {
		t.Skip("no sound on this machine:", err)
	}
	if err != nil {
		t.Fatal(err)
	}

	queuesMu.Lock()
	open := len(queues)
	queuesMu.Unlock()
	if open != 1 {
		t.Fatalf("%d streams are open, want one", open)
	}

	if err := audio.Close(); err != nil {
		t.Fatal(err)
	}
	queuesMu.Lock()
	open = len(queues)
	queuesMu.Unlock()
	if open != 0 {
		t.Errorf("%d streams survived closing", open)
	}

	// A late callback for a stream that has gone silences its buffer instead
	// of following a nil pointer into the mixer.
	buffer := make([]float32, 8)
	for i := range buffer {
		buffer[i] = 0.5
	}
	fillQueue(9999, buffer)
	for i, sample := range buffer {
		if sample != 0 {
			t.Errorf("sample %d is %v, want silence", i, sample)
		}
	}
}