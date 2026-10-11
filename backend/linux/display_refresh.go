//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
)

// DisplayRefresh is how often the display repaints, in hertz. 0 when it
// cannot be known, which is any server without RandR.
func (x *Driver) DisplayRefresh() int {
	if !x.alive {
		return 0
	}
	if x.refreshHz != 0 {
		return max(x.refreshHz, 0)
	}
	x.refreshHz = -1 // asked, and unknown unless RandR says otherwise
	if x.randr() == 0 {
		return 0
	}
	seq := x.send(x.randrOpcode, xRRGetScreenInfo, le32(x.root))
	reply, _, ok := x.awaitReply(seq)
	if !ok {
		return 0
	}
	if rate := binary.LittleEndian.Uint16(reply[26:]); rate > 0 && rate < 1000 {
		x.refreshHz = int(rate)
	}
	return max(x.refreshHz, 0)
}
