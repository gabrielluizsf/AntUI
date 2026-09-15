package audio

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// WAV is the container the sounds ship in: a tiny RIFF header and then plain
// PCM samples. Decoding it is straight reading, and the files were converted
// once (from the source Ogg Vorbis) so nothing else ever has to be.

// decodeWAV turns a RIFF/WAVE file into mono float samples in [-1, 1] plus
// its sample rate. Any PCM layout the encoder may have left — mono or stereo,
// 8 or 16 bit — is accepted and mixed down to mono.
func decodeWAV(data []byte) (samples []float32, rate int, err error) {
	if len(data) < 12 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WAVE" {
		return nil, 0, errors.New("audio: not a RIFF/WAVE file")
	}

	var channels, bits uint16
	haveFmt := false

	for pos := 12; pos+8 <= len(data); {
		chunk := data[pos : pos+4]
		size := int(binary.LittleEndian.Uint32(data[pos+4 : pos+8]))
		body := pos + 8
		if body+size > len(data) {
			return nil, 0, errors.New("audio: truncated WAVE file")
		}
		switch string(chunk) {
		case "fmt ":
			if size < 16 {
				return nil, 0, errors.New("audio: short fmt chunk")
			}
			// Only PCM is ever shipped here. A compressed format has a
			// block-aligned payload this decoder does not understand, and
			// refusing is honest where guessing would sound wrong.
			if format := binary.LittleEndian.Uint16(data[body : body+2]); format != 1 {
				return nil, 0, fmt.Errorf("audio: compressed WAVE format %d", format)
			}
			channels = binary.LittleEndian.Uint16(data[body+2 : body+4])
			rate = int(binary.LittleEndian.Uint32(data[body+4 : body+8]))
			bits = binary.LittleEndian.Uint16(data[body+14 : body+16])
			if channels == 0 || channels > 2 {
				return nil, 0, fmt.Errorf("audio: %d channels", channels)
			}
			haveFmt = true
		case "data":
			if !haveFmt {
				return nil, 0, errors.New("audio: samples before the format")
			}
			raw := data[body : body+size]
			stride := int(channels) * int(bits/8)
			if stride == 0 || len(raw)%stride != 0 {
				return nil, 0, errors.New("audio: misaligned sample data")
			}
			samples = make([]float32, 0, len(raw)/stride)
			switch bits {
			case 8:
				for i := 0; i+int(channels) <= len(raw); i += int(channels) {
					sum := 0
					for c := range channels {
						sum += int(int16(raw[i+int(c)]) - 128) // signed, 8-bit is stored unsigned
					}
					samples = append(samples, float32(sum)/float32(channels)/128)
				}
			case 16:
				for i := 0; i+stride <= len(raw); i += stride {
					sum := 0
					for c := range channels {
						sum += int(int16(binary.LittleEndian.Uint16(raw[i+2*int(c):])))
					}
					samples = append(samples, float32(sum)/float32(channels)/32768)
				}
			default:
				return nil, 0, fmt.Errorf("audio: %d-bit samples", bits)
			}
			normalize(samples)
			return samples, rate, nil
		}
		pos = body + size
		if size%2 != 0 {
			pos++ // chunks are word-aligned
		}
	}
	return nil, 0, errors.New("audio: WAVE file with no samples")
}

// normalize lifts a quiet recording up to a full swing, so a bank's sounds
// all come out at the same presence whatever the source level was.
func normalize(samples []float32) {
	var peak float32
	for _, s := range samples {
		if a := abs(s); a > peak {
			peak = a
		}
	}
	if peak == 0 || peak >= 0.99 {
		return
	}
	gain := 0.92 / peak
	for i := range samples {
		samples[i] *= gain
	}
}

func abs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}