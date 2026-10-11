//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// detectPixels walks the screen's visuals to find the root visual's channel
// masks, then works out whether the server's bytes are already what the
// canvas holds.
func (x *Driver) detectPixels(body []byte, screenAt int, rootVisual uint32) {
	// The root visual's channel masks say whether red and blue sit the other
	// way round from the order the canvas holds them in.
	var redMask, blueMask uint32
	at := screenAt + 40
	for range int(body[screenAt+39]) {
		if at+8 > len(body) {
			break
		}
		visuals := int(binary.LittleEndian.Uint16(body[at+2:]))
		at += 8
		for range visuals {
			if at+24 > len(body) {
				break
			}
			if binary.LittleEndian.Uint32(body[at:]) == rootVisual {
				redMask = binary.LittleEndian.Uint32(body[at+8:])
				blueMask = binary.LittleEndian.Uint32(body[at+16:])
			}
			at += 24
		}
	}
	x.swapRB = redMask == 0x000000FF && blueMask == 0x00FF0000
	x.samePixels = x.imageByteOrder == 0 && !x.swapRB &&
		binary.NativeEndian.Uint32([]byte{1, 0, 0, 0}) == 1
}
