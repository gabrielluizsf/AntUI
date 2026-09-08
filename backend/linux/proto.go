//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"time"
)

// awaitReply waits for the answer to one request, handling any events that
// arrive ahead of it rather than dropping them on the floor.
func (x *Driver) awaitReply(seq uint16) (reply, extra []byte, ok bool) {
	deadline := time.Now().Add(xReplyTimeout)
	for {
		packet, body, got := x.takePacket()
		if !got {
			if !x.fill(time.Until(deadline)) {
				if !x.alive || !time.Now().Before(deadline) {
					return nil, nil, false
				}
			}
			continue
		}
		packetSeq := binary.LittleEndian.Uint16(packet[2:])
		switch {
		case packet[0] == xReply && packetSeq == seq:
			return packet, body, true
		case packet[0] == xError && packetSeq == seq:
			return nil, nil, false
		case packet[0] > xReply && x.win != nil:
			x.handleEvent(packet)
		}
	}
}

// request builds a request whose length field it fills in itself. The body
// is padded to a four-byte boundary, as every request must be.
func (x *Driver) request(opcode, detail byte, body []byte) []byte {
	total := 4 + len(body)
	total += pad(total)
	buf := make([]byte, total)
	buf[0] = opcode
	buf[1] = detail
	binary.LittleEndian.PutUint16(buf[2:], uint16(total/4))
	copy(buf[4:], body)
	return buf
}

// send writes a request and returns the sequence number the server will
// answer it under.
func (x *Driver) send(opcode, detail byte, body []byte) uint16 {
	x.seq++
	x.write(x.request(opcode, detail, body))
	return x.seq
}

func le32(values ...uint32) []byte {
	buf := make([]byte, 4*len(values))
	for i, v := range values {
		binary.LittleEndian.PutUint32(buf[i*4:], v)
	}
	return buf
}

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

// changeProperty sets a property on the window. count is in items of the
// size named by format, not in bytes.
func (x *Driver) changeProperty(property, typ uint32, format byte, value []byte, count uint32) {
	x.changePropertyOn(x.window, property, typ, format, value, count)
}

// changePropertyOn is the same on somebody else's window, which is how a
// selection is handed over: the asker names a property on their own window
// and the owner writes the answer into it.
func (x *Driver) changePropertyOn(window, property, typ uint32, format byte,
	value []byte, count uint32) {

	body := make([]byte, 20, 20+len(value)+3)
	binary.LittleEndian.PutUint32(body[0:], window)
	binary.LittleEndian.PutUint32(body[4:], property)
	binary.LittleEndian.PutUint32(body[8:], typ)
	body[12] = format
	binary.LittleEndian.PutUint32(body[16:], count)
	body = append(body, value...)
	x.send(xChangeProperty, 0, body) // detail 0 = replace
}

// setTitleOn writes the title twice: WM_NAME for window managers that only
// read Latin-1, and _NET_WM_NAME in UTF-8 for everything from this century.
func (x *Driver) setTitleOn(title string) {
	latin1 := make([]byte, 0, len(title))
	for _, r := range title {
		if r <= 0xFF {
			latin1 = append(latin1, byte(r))
		} else {
			latin1 = append(latin1, '?')
		}
	}
	x.changeProperty(atomWMName, atomString, 8, latin1, uint32(len(latin1)))
	if x.atomNetWMName != 0 && x.atomUTF8String != 0 {
		x.changeProperty(x.atomNetWMName, x.atomUTF8String, 8, []byte(title), uint32(len(title)))
	}
}

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