//go:build windows

package antui

import (
	"errors"
	"runtime"
	"testing"
	"time"
	"unsafe"
)

// These run on Windows and nowhere else. The rest of the audio tests are the
// same on every platform; what is here is the part that can only be wrong in
// a WASAPI backend — the sample layouts, and whether the COM handshake it is
// written from is the one Windows actually has.

// The mix format is whatever the shared engine is running at, and the two
// layouts it is ever in are 32-bit floats and 16-bit integers.
func TestWASAPIUnderstandsTheMixFormats(t *testing.T) {
	w := &wasapi{}
	for _, c := range []struct {
		name   string
		format waveFormatExtensible
		known  bool
		floats bool
	}{
		{"plain floats", waveFormatExtensible{waveFormatEx: waveFormatEx{
			tag: formatFloat, channels: 2, samplesPerSec: 48000, bitsPerSample: 32}},
			true, true},
		{"plain 16-bit", waveFormatExtensible{waveFormatEx: waveFormatEx{
			tag: formatPCM, channels: 2, samplesPerSec: 44100, bitsPerSample: 16}},
			true, false},
		{"extensible floats", waveFormatExtensible{waveFormatEx: waveFormatEx{
			tag: formatExtensible, channels: 6, samplesPerSec: 48000, bitsPerSample: 32},
			subFormat: subtypeFloat}, true, true},
		{"extensible 16-bit", waveFormatExtensible{waveFormatEx: waveFormatEx{
			tag: formatExtensible, channels: 2, samplesPerSec: 48000, bitsPerSample: 16},
			subFormat: subtypePCM}, true, false},
		{"24-bit", waveFormatExtensible{waveFormatEx: waveFormatEx{
			tag: formatExtensible, channels: 2, samplesPerSec: 48000, bitsPerSample: 24},
			subFormat: subtypePCM}, false, false},
		{"nothing at all", waveFormatExtensible{}, false, false},
	} {
		format := c.format
		if got := w.understands(&format); got != c.known {
			t.Errorf("%s: understood %v, want %v", c.name, got, c.known)
		}
		if c.known {
			if got := w.floats(&format); got != c.floats {
				t.Errorf("%s: floats %v, want %v", c.name, got, c.floats)
			}
		}
	}
}

// A stereo mix goes to the front pair of however many channels the device
// has, and the rest are left silent: the same sound in the surrounds as well
// makes a room sound like a corridor.
func TestWASAPISpreadsStereoAcrossTheDevice(t *testing.T) {
	const frames = 4
	block := make([]float32, frames*2)
	for i := range block {
		block[i] = float32(i+1) / 16 // 0.0625, 0.125, ...
	}

	for _, channels := range []int{1, 2, 6} {
		buffer := make([]float32, frames*channels)
		for i := range buffer {
			buffer[i] = 99 // what a buffer holds when it arrives is anyone's
		}
		w := &wasapi{}
		w.spread(uintptr(unsafe.Pointer(&buffer[0])), block, frames, channels, true)
		runtime.KeepAlive(buffer)

		for frame := range frames {
			if got, want := buffer[frame*channels], block[frame*2]; got != want {
				t.Errorf("%d channels: frame %d left is %v, want %v", channels, frame, got, want)
			}
			if channels > 1 {
				if got, want := buffer[frame*channels+1], block[frame*2+1]; got != want {
					t.Errorf("%d channels: frame %d right is %v, want %v", channels, frame, got, want)
				}
			}
			for c := 2; c < channels; c++ {
				if got := buffer[frame*channels+c]; got != 0 {
					t.Errorf("%d channels: frame %d channel %d is %v, want silence",
						channels, frame, c, got)
				}
			}
		}
	}
}

// The 16-bit path clips rather than wrapping. A sample that overflows into
// the opposite sign is a click, and a click is the loudest thing a speaker
// can make.
func TestWASAPIClipsRatherThanWraps(t *testing.T) {
	for _, c := range []struct {
		sample float32
		want   int16
	}{
		{0, 0}, {1, 32767}, {-1, -32767}, {2, 32767}, {-2, -32768},
		{0.5, 16384}, {-0.5, -16384},
	} {
		if got := pcm16(c.sample); got != c.want {
			t.Errorf("%v came out as %d, want %d", c.sample, got, c.want)
		}
	}
}

// And the whole of it against the machine's own device: the COM vtable
// indexes this is written from are either right or nothing plays, and there
// is no way to find out but to ask Windows.
func TestWASAPIPlaysOnThisMachine(t *testing.T) {
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

	t.Logf("the device is at %d Hz with %v of latency", audio.Rate(), audio.Latency())
	if audio.Rate() < 8000 || audio.Rate() > 384000 {
		t.Errorf("the device came back at %d Hz", audio.Rate())
	}
	// Shared mode and event driven is about ten milliseconds. Anything near
	// waveOut's forty means the event path is not the one being taken.
	if got := audio.Latency(); got <= 0 || got > 50*time.Millisecond {
		t.Errorf("the latency is %v, which is not what shared mode gives", got)
	}

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatalf("the device asked for %d frames in three seconds", frames)
	}
	if err := audio.Err(); err != nil {
		t.Fatal(err)
	}
}