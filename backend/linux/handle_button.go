//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"

	"github.com/gabrielluizsf/antui/backend"
)

// handleButton folds a button press or release into the window. The wheel
// arrives as buttons 4 and 5, and only on press.
func (x *Driver) handleButton(packet []byte) {
	press := packet[0]&0x7F == xButtonPress
	state := binary.LittleEndian.Uint16(packet[28:])
	ev := backend.Event{
		Mods: modsFromState(state),
		X:    int(int16(binary.LittleEndian.Uint16(packet[24:]))),
		Y:    int(int16(binary.LittleEndian.Uint16(packet[26:]))),
	}
	switch button := packet[1]; button {
	case 4, 5:
		if press {
			ev.Type = backend.EventMouseWheel
			ev.Wheel = 1
			if button == 5 {
				ev.Wheel = -1
			}
			x.win.Push(ev)
		}
		return
	case 1:
		ev.Button = backend.MouseLeft
	case 2:
		ev.Button = backend.MouseMiddle
	case 3:
		ev.Button = backend.MouseRight
	default:
		return
	}
	ev.Type = backend.EventMouseUp
	if press {
		ev.Type = backend.EventMouseDown
	}
	x.win.Push(ev)
}
