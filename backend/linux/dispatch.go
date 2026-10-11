//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"github.com/gabrielluizsf/antui/backend"
)

// dispatch folds every whole packet in the buffer into the window. A packet
// half-read is left for the next time: the reply it belongs to may still be
// coming, and a read cannot be allowed to wait for one here.
func (x *Driver) dispatch(win backend.Face) {
	for {
		packet, _, ok := x.takePacket()
		if !ok {
			// Half a packet is buffered: try to complete it without blocking.
			if !x.fill(0) {
				return
			}
			continue
		}
		if packet[0] != xReply {
			x.handleEvent(packet)
		}
	}
}
