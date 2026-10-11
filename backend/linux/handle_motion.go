//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"

	"github.com/gabrielluizsf/antui/backend"
)

// handleMotion folds a pointer move into the window.
func (x *Driver) handleMotion(packet []byte) {
	x.win.Push(backend.Event{
		Type: backend.EventMouseMove,
		X:    int(int16(binary.LittleEndian.Uint16(packet[24:]))),
		Y:    int(int16(binary.LittleEndian.Uint16(packet[26:]))),
		Mods: modsFromState(binary.LittleEndian.Uint16(packet[28:])),
	})
}
