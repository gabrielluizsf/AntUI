//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// loadKeymap fetches the keycode-to-keysym table. It is read once at startup
// and again whenever the server reports the layout changed, so that switching
// keyboard layout mid-run does what the user expects.
func (x *Driver) loadKeymap() bool {
	count := int(x.maxKeycode) - int(x.minKeycode) + 1
	if count <= 0 {
		return false
	}
	seq := x.send(xGetKeyboardMap, 0, []byte{x.minKeycode, byte(count), 0, 0})
	reply, extra, ok := x.awaitReply(seq)
	if !ok {
		return false
	}
	x.keysymsPerCode = reply[1]
	x.keysyms = nil

	want := count * int(x.keysymsPerCode)
	if x.keysymsPerCode == 0 || len(extra) < want*4 {
		return false
	}
	x.keysyms = make([]uint32, want)
	for i := range x.keysyms {
		x.keysyms[i] = binary.LittleEndian.Uint32(extra[i*4:])
	}
	return true
}
