package antui

import (
	"errors"
	"math"
	"sync/atomic"
	"testing"
	"time"
)

// A device that opens has to actually ask for sound, and go on asking. Half
// of what can go wrong with a backend — a handshake that half worked, a
// buffer the server never asks to have filled — looks exactly like a working
// one until you count the frames it took.
//
// Skipped where there is no sound at all, which is most build machines: a
// test that fails on a box with no speakers is a test nobody can run.
func TestAudioAsksForSound(t *testing.T) {
	const rate = 44100

	var frames atomic.Int64
	phase := 0.0
	audio, err := OpenAudio("AntUI test", rate, func(block []float32) {
		// A quarter-volume tone at 440, so that anyone who does run this
		// with the speakers on hears something recognisable rather than
		// noise.
		for i := 0; i < len(block); i += 2 {
			sample := float32(math.Sin(phase) * 0.25)
			block[i], block[i+1] = sample, sample
			phase += 2 * math.Pi * 440 / rate
		}
		frames.Add(int64(len(block) / 2))
	})
	if errors.Is(err, ErrNoAudio) {
		t.Skip("no sound on this machine:", err)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer audio.Close()

	if got := audio.Rate(); got != rate {
		t.Errorf("the device is running at %d, want %d", got, rate)
	}
	if got := audio.Channels(); got != 2 {
		t.Errorf("the device has %d channels, want stereo", got)
	}
	if got := audio.Latency(); got <= 0 || got > 200*time.Millisecond {
		t.Errorf("the latency is %v, which is not a latency a game can use", got)
	}

	// A tenth of a second of sound, which a device that is really running
	// hands over in about a tenth of a second.
	deadline := time.Now().Add(3 * time.Second)
	for frames.Load() < rate/10 {
		if time.Now().After(deadline) {
			t.Fatalf("the device asked for %d frames in three seconds", frames.Load())
		}
		if err := audio.Err(); err != nil {
			t.Fatal(err)
		}
		time.Sleep(5 * time.Millisecond)
	}
	if err := audio.Err(); err != nil {
		t.Fatal(err)
	}

	// And it stops when it is told to, without the goroutine behind it
	// deciding that a closed socket is a disaster worth reporting.
	if err := audio.Close(); err != nil {
		t.Errorf("closing: %v", err)
	}
	time.Sleep(50 * time.Millisecond)
	if err := audio.Err(); err != nil {
		t.Errorf("closing the device reported %v", err)
	}
	if err := audio.Close(); err != nil {
		t.Errorf("closing twice: %v", err)
	}
}

// Opening with nothing to play is a mistake worth an error rather than a
// device that pulls silence for ever.
func TestAudioNeedsSomethingToPlay(t *testing.T) {
	if _, err := OpenAudio("AntUI test", 44100, nil); err == nil {
		t.Error("a device opened with no way to fill it")
	}
}