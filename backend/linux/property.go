//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
)

// property reads a property off a window. maxBytes bounds what is asked for
// in one go: a property longer than that is read short rather than in
// several round trips, which is right for everything here — the resource
// database is the only one read at all, and 16 KB of it is far more than the
// one line this wants.
func (x *Driver) property(window, name, typ uint32, maxBytes int) ([]byte, bool) {
	if window == 0 || name == 0 {
		return nil, false
	}
	body := make([]byte, 20)
	binary.LittleEndian.PutUint32(body[0:], window)
	binary.LittleEndian.PutUint32(body[4:], name)
	binary.LittleEndian.PutUint32(body[8:], typ)
	binary.LittleEndian.PutUint32(body[12:], 0) // offset, in 32-bit words
	binary.LittleEndian.PutUint32(body[16:], uint32(maxBytes/4))

	seq := x.send(xGetProperty, 0, body) // detail 0 = do not delete it
	reply, extra, ok := x.awaitReply(seq)
	if !ok || binary.LittleEndian.Uint32(reply[8:]) == 0 {
		return nil, false // no such property, or not of the type asked for
	}
	length := int(binary.LittleEndian.Uint32(reply[16:])) * int(reply[1]) / 8
	if length <= 0 || length > len(extra) {
		return nil, false
	}
	return extra[:length], true
}
