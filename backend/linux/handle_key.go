//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"

	"github.com/gabrielluizsf/antui/backend"
)

// handleKey turns one key packet into a key event and, when the key stands
// for a character, the text event that goes with it.
func (x *Driver) handleKey(packet []byte) {
	win := x.win
	if win == nil {
		return
	}
	press := packet[0]&0x7F == xKeyPress
	state := binary.LittleEndian.Uint16(packet[28:])
	sym := x.keysym(packet[1], state)

	ev := backend.Event{
		Type: backend.EventKeyUp,
		Key:  keysymToKey(sym),
		Mods: modsFromState(state),
		X:    int(int16(binary.LittleEndian.Uint16(packet[24:]))),
		Y:    int(int16(binary.LittleEndian.Uint16(packet[26:]))),
	}
	if press {
		ev.Type = backend.EventKeyDown
	}
	win.Push(ev)
	if !press {
		return
	}

	// A dead accent produces no text of its own: it waits for the next key.
	if sym >= symDeadLo && sym <= symDeadHi {
		x.deadKey = sym
		return
	}
	x.pushKeyText(win, sym, ev.Mods)
}
