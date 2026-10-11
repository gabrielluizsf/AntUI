//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
)

// readProperty reads a property off our own window and deletes it, which is
// what a selection's protocol asks of whoever receives one.
func (x *Driver) readProperty(property uint32) []byte {
	var out []byte
	for offset := uint32(0); ; {
		// Delete on the last read only, so a property that needs two reads is
		// not thrown away halfway through.
		seq := x.send(xGetProperty, 0,
			le32(x.window, property, 0, offset, 4096))
		reply, extra, ok := x.awaitReply(seq)
		if !ok {
			break
		}
		// The reply pads its data to a four-byte boundary, and a property
		// whose length is not a multiple of four comes back with the padding
		// still attached — a 26-character copy as 28 bytes with two zeros on
		// the end. The reply says how many items are really there, which is
		// what is cut to. (The original port of this function carried the
		// padding; the integration test caught it on a read-back.)
		items := int(binary.LittleEndian.Uint32(reply[16:]))
		length := items * int(reply[1]) / 8
		if length < 0 || length > len(extra) {
			break
		}
		out = append(out, extra[:length]...)
		if binary.LittleEndian.Uint32(reply[12:]) == 0 { // bytes-after
			break
		}
		offset += uint32(length / 4)
	}
	// Deleting says "I have taken it", which is what the holder waits for.
	x.send(xGetProperty, 1, le32(x.window, property, 0, 0, 0))
	return out
}
