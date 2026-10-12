package audio

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// readSamples decodes a data chunk into mono float samples, mixing the
// channels down and scaling by the bit depth.
func readSamples(raw []byte, channels, bits uint16) ([]float32, error) {
	stride := int(channels) * int(bits/8)
	if stride == 0 || len(raw)%stride != 0 {
		return nil, errors.New("audio: misaligned sample data")
	}
	samples := make([]float32, 0, len(raw)/stride)
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
		return nil, fmt.Errorf("audio: %d-bit samples", bits)
	}
	return samples, nil
}
