//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"time"
)

// askSelection asks whoever holds a selection to write it onto a property of
// our window, waits for them to say they have, and reads it.
//
// The wait is bounded, and generously: a selection is held by another process,
// and a process that has stopped answering must not stop this one — but a
// busy machine can take a while to schedule it, and a drop that quietly
// brings in nothing because the answer arrived a moment late is worse than a
// frame that took a second. A second is long past what any live source needs
// and far short of a hang.
func (x *Driver) askSelection(selection, target, when uint32) []byte {
	if !x.alive || selection == 0 || target == 0 {
		return nil
	}
	x.send(xConvertSelection, 0, le32(x.window, selection, target,
		x.atomDropProperty, when))

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		packet, _, ok := x.takePacket()
		if !ok {
			x.fill(20 * time.Millisecond)
			continue
		}
		if packet[0]&0x7F != xSelectionNotify {
			if packet[0] != xReply && x.win != nil {
				x.handleEvent(packet)
			}
			continue
		}
		// The property is None when the holder could not give us what was
		// asked for — a file manager with text on the clipboard and nothing
		// that looks like a file, say.
		if binary.LittleEndian.Uint32(packet[20:]) == 0 {
			return nil
		}
		return x.readProperty(x.atomDropProperty)
	}
	return nil
}
