//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
)

// internAtom asks the server for the number it knows a name by. Atoms are
// how everything outside the core protocol — window manager hints, the UTF-8
// title, full screen — is addressed.
func (x *Driver) internAtom(name string) uint32 {
	body := make([]byte, 4, 4+len(name)+3)
	binary.LittleEndian.PutUint16(body, uint16(len(name)))
	body = append(body, name...)

	seq := x.send(xInternAtom, 0, body) // detail 0 = create it if it is new
	reply, _, ok := x.awaitReply(seq)
	if !ok {
		return 0
	}
	return binary.LittleEndian.Uint32(reply[8:])
}
