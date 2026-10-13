//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

func (x *Driver) sendDragMessage(kind uint32, data [5]uint32) {
	body := make([]byte, xClientMessageBytes)
	binary.LittleEndian.PutUint32(body[0:], x.dragFrom)
	binary.LittleEndian.PutUint32(body[4:], 0) // straight to the client

	event := body[xClientMessageHeaderBytes:]
	event[0] = xClientMessage
	event[1] = xFormat32Bits // format
	binary.LittleEndian.PutUint32(event[4:], x.dragFrom)
	binary.LittleEndian.PutUint32(event[8:], kind)
	for i, value := range data {
		binary.LittleEndian.PutUint32(event[12+i*4:], value)
	}
	x.send(xSendEvent, 0, body)
}
