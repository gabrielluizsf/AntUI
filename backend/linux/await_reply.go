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
