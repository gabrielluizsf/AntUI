package audio

import (
	"math"
	"testing"
)

func pcm16LE(samples []int16) []byte {
	data := make([]byte, len(samples)*2)
	for i, s := range samples {
		data[i*2] = byte(s)
		data[i*2+1] = byte(s >> 8)
	}
	return data
}

func monoWav16(rate int, samples []int16) []byte {
	raw := pcm16LE(samples)
	channels := 1
	bits := 16
	// RIFF header 12 + fmt chunk 24 + data chunk header 8 + data
	size := len(raw)
	total := 12 + 24 + 8 + size
	out := make([]byte, total)
	copy(out[0:4], "RIFF")
	putU32LE(out[4:], uint32(total-8))
	copy(out[8:12], "WAVE")
	// fmt chunk
	copy(out[12:16], "fmt ")
	putU32LE(out[16:], 16)
	putU16LE(out[20:], 1) // PCM
	putU16LE(out[22:], uint16(channels))
	putU32LE(out[24:], uint32(rate))
	putU32LE(out[28:], uint32(rate*channels*bits/8))
	putU16LE(out[32:], uint16(channels*bits/8))
	putU16LE(out[34:], uint16(bits))
	// data chunk
	copy(out[36:40], "data")
	putU32LE(out[40:], uint32(size))
	copy(out[44:], raw)
	return out
}

func putU16LE(b []byte, v uint16) { b[0], b[1] = byte(v), byte(v>>8) }
func putU32LE(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func TestDecodeMono16(t *testing.T) {
	src := []int16{0, 32767, -32767, 0}
	data := monoWav16(22050, src)
	samples, rate, err := decodeWAV(data)
	if err != nil {
		t.Fatal(err)
	}
	if rate != 22050 {
		t.Errorf("rate = %d want 22050", rate)
	}
	if len(samples) != len(src) {
		t.Fatalf("len(samples) = %d want %d", len(samples), len(src))
	}
	// Full-scale samples: peak is 1.0, so normalisation is a no-op.
	if math.Abs(float64(samples[1]-1.0)) > 0.01 {
		t.Errorf("samples[1] = %f want ~1.0", samples[1])
	}
	if math.Abs(float64(samples[2]+1.0)) > 0.01 {
		t.Errorf("samples[2] = %f want ~-1.0", samples[2])
	}
}

func TestDecodeRejectsRIFFMagic(t *testing.T) {
	_, _, err := decodeWAV([]byte("NOT A WAV FILE"))
	if err == nil {
		t.Error("expected error for missing RIFF header")
	}
}

func TestDecodeRejectsCompressedFormat(t *testing.T) {
	data := monoWav16(44100, []int16{0})
	// patch fmt codec to 0xFFFE (extensible)
	data[20] = 0xFE
	data[21] = 0xFF
	_, _, err := decodeWAV(data)
	if err == nil {
		t.Error("expected error for non-PCM format")
	}
}

func TestNormalizeLiftsQuietSound(t *testing.T) {
	data := monoWav16(44100, []int16{0, 100, 200, 100})
	samples, _, err := decodeWAV(data)
	if err != nil {
		t.Fatal(err)
	}
	var peak float32
	for _, s := range samples {
		if a := s; a > peak {
			peak = a
		}
	}
	// normalize() should have pushed the quietest sample close to the 0.92 target.
	if peak < 0.5 {
		t.Errorf("normalization did not lift the peak, got %f", peak)
	}
}

func TestDecodeHandlesStereo16(t *testing.T) {
	// Stereo WAV: L=32767, R=-32767 → mono = 0.
	raw := []byte{0xFF, 0x7F, 0x01, 0x80}
	channels, bits, rate := 2, 16, 44100
	data := make([]byte, 44+len(raw))
	copy(data[0:4], "RIFF")
	putU32LE(data[4:], uint32(44-8+len(raw)))
	copy(data[8:12], "WAVE")
	copy(data[12:16], "fmt ")
	putU32LE(data[16:], 16)
	putU16LE(data[20:], 1)
	putU16LE(data[22:], uint16(channels))
	putU32LE(data[24:], uint32(rate))
	putU32LE(data[28:], uint32(rate*channels*bits/8))
	putU16LE(data[32:], uint16(channels*bits/8))
	putU16LE(data[34:], uint16(bits))
	copy(data[36:40], "data")
	putU32LE(data[40:], uint32(len(raw)))
	copy(data[44:], raw)
	samples, _, err := decodeWAV(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(samples) != 1 {
		t.Fatalf("expected 1 mono frame, got %d", len(samples))
	}
	// The stereo mix of +1 and -1 should cancel.
	if math.Abs(float64(samples[0])) > 0.01 {
		t.Errorf("expected near-zero, got %f", samples[0])
	}
}