package font

import "fmt"

// ParseTTF reads the tables a face needs out of a font file.
func ParseTTF(data []byte) (*TTF, error) {
	if len(data) < 12 {
		return nil, fmt.Errorf("antUI: %d bytes is not a font", len(data))
	}
	switch tag := be32(data, 0); tag {
	case 0x00010000, 0x74727565: // 1.0, and "true"
	case 0x74746366: // "ttcf": a collection, whose first face is taken
		if len(data) < 16 {
			return nil, fmt.Errorf("antUI: a truncated font collection")
		}
		offset := be32(data, 12)
		if int(offset) >= len(data) {
			return nil, fmt.Errorf("antUI: a font collection pointing past its end")
		}
		return parseTTFAt(data, int(offset))
	case 0x4F54544F: // "OTTO"
		return nil, fmt.Errorf("antUI: this is a CFF font, and only TrueType " +
			"outlines are read here")
	default:
		return nil, fmt.Errorf("antUI: %#08x is not a font's tag", tag)
	}
	return parseTTFAt(data, 0)
}
