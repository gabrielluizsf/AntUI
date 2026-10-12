package audio

import (
	"encoding/binary"
	"errors"
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
			if channels, bits, rate, err = readFmtChunk(data[body : body+size]); err != nil {
				return nil, 0, err
			}
			haveFmt = true
		case "data":
			if !haveFmt {
				return nil, 0, errors.New("audio: samples before the format")
			}
			samples, err = readSamples(data[body:body+size], channels, bits)
			if err != nil {
				return nil, 0, err
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
