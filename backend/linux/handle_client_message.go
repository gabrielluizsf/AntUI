//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"

	"github.com/gabrielluizsf/antui/backend"
)

// handleClientMessage answers a client message: a request to close, or a
// drag event for this window to react to.
func (x *Driver) handleClientMessage(packet []byte) {
	typeAtom := binary.LittleEndian.Uint32(packet[8:])
	data0 := binary.LittleEndian.Uint32(packet[12:])
	if typeAtom == x.atomWMProtocols && data0 == x.atomWMDeleteWindow {
		x.win.PushSimple(backend.EventClose)
		return
	}
	x.handleDrop(packet)
}
