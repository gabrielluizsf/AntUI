//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
)

// sendSelectionNotify tells the asker the answer is on their window, or that
// there is no answer.
func (x *Driver) sendSelectionNotify(requestor, property, target, when uint32) {
	body := make([]byte, 40)
	binary.LittleEndian.PutUint32(body[0:], requestor)
	binary.LittleEndian.PutUint32(body[4:], 0) // straight to the client

	event := body[8:]
	event[0] = xSelectionNotify
	binary.LittleEndian.PutUint32(event[4:], when)
	binary.LittleEndian.PutUint32(event[8:], requestor)
	binary.LittleEndian.PutUint32(event[12:], x.atomClipboard)
	binary.LittleEndian.PutUint32(event[16:], target)
	binary.LittleEndian.PutUint32(event[20:], property)

	x.send(xSendEvent, 0, body)
}
