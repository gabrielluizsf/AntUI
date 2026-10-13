//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// wmState asks the window manager to turn a state on or off.
//
// EWMH says a mapped window changes state by asking the window manager,
// through a ClientMessage sent to the root. Setting the property directly
// only works before the window is mapped, and this window is mapped by the
// time anybody asks.
func (x *Driver) wmState(action, state uint32) bool {
	if !x.alive || x.atomNetWMState == 0 || state == 0 {
		return false
	}
	body := make([]byte, 40)
	binary.LittleEndian.PutUint32(body[0:], x.root)
	// SubstructureNotify | SubstructureRedirect: what a root-directed EWMH
	// message is required to be sent with.
	binary.LittleEndian.PutUint32(body[4:], 0x180000)

	event := body[8:]
	event[0] = xClientMessage
	event[1] = 32 // format
	binary.LittleEndian.PutUint32(event[4:], x.window)
	binary.LittleEndian.PutUint32(event[8:], x.atomNetWMState)
	binary.LittleEndian.PutUint32(event[12:], action) // 0 off, 1 on, 2 flip
	binary.LittleEndian.PutUint32(event[16:], state)
	binary.LittleEndian.PutUint32(event[24:], 1) // source = application

	x.send(xSendEvent, 0, body) // detail 0 = do not propagate
	return x.alive
}
