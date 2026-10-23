package svg

import "errors"

// gifFrameCount walks the blocks of a GIF without decoding any of them and
// answers how many frames it holds. Every frame becomes a whole picture the
// size of the GIF's screen once it is read, so this is what the size of the
// whole animation is worked out from before a byte of it is decompressed.
func gifFrameCount(data []byte) (int, error) {
	errShort := errors.New("it is a GIF that ends early")
	if len(data) < 13 {
		return 0, errShort
	}
	p := 13
	if flags := data[10]; flags&0x80 != 0 {
		p += 3 << (flags&7 + 1) // the global colour table
	}
	frames := 0
	for p < len(data) {
		switch data[p] {
		case 0x3B: // the trailer
			return frames, nil
		case 0x21: // an extension: its label, then its sub-blocks
			p += 2
		case 0x2C: // an image: its descriptor, its own colour table, the code size
			if p+10 > len(data) {
				return 0, errShort
			}
			frames++
			flags := data[p+9]
			p += 10
			if flags&0x80 != 0 {
				p += 3 << (flags&7 + 1)
			}
			p++
		default:
			return 0, errors.New("it is a GIF with a block it does not name")
		}
		for p < len(data) && data[p] != 0 {
			p += int(data[p]) + 1
		}
		p++
	}
	return frames, nil
}
