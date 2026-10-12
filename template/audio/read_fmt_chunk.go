package audio

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// readFmtChunk reads the channel count, bit depth and sample rate out of one
// fmt chunk.
func readFmtChunk(body []byte) (channels, bits uint16, rate int, err error) {
	if len(body) < 16 {
		return 0, 0, 0, errors.New("audio: short fmt chunk")
	}
	// Only PCM is ever shipped here. A compressed format has a
	// block-aligned payload this decoder does not understand, and
	// refusing is honest where guessing would sound wrong.
	if format := binary.LittleEndian.Uint16(body[0:2]); format != 1 {
		return 0, 0, 0, fmt.Errorf("audio: compressed WAVE format %d", format)
	}
	channels = binary.LittleEndian.Uint16(body[2:4])
	rate = int(binary.LittleEndian.Uint32(body[4:8]))
	bits = binary.LittleEndian.Uint16(body[14:16])
	if channels == 0 || channels > 2 {
		return 0, 0, 0, fmt.Errorf("audio: %d channels", channels)
	}
	return channels, bits, rate, nil
}
