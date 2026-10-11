//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
)

// randr looks up the RandR extension's opcode, once. The core protocol has
// no refresh rate, so without the extension the honest answer is that we do
// not know.
func (x *Driver) randr() byte {
	if x.randrOpcode != 0 {
		return x.randrOpcode
	}
	const name = "RANDR"
	body := make([]byte, 4, 12)
	binary.LittleEndian.PutUint16(body, uint16(len(name)))
	body = append(body, name...)

	seq := x.send(xQueryExtension, 0, body)
	reply, _, ok := x.awaitReply(seq)
	if !ok || reply[8] == 0 { // present = false
		return 0
	}
	x.randrOpcode = reply[9]

	// Some servers will not answer anything else until a version has been
	// agreed on. 1.1 is enough for a refresh rate and is twenty years old.
	seq = x.send(x.randrOpcode, xRRQueryVersion, le32(1, 1))
	x.awaitReply(seq)
	return x.randrOpcode
}
