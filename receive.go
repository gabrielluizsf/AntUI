//go:build linux && !android

package antui

import (
	"encoding/binary"
	"fmt"
	"io"
)

// receive reads one packet.
func (p *pulse) receive() (channel uint32, payload []byte, err error) {
	var head [20]byte
	if _, err := io.ReadFull(p.conn, head[:]); err != nil {
		return 0, nil, err
	}
	length := binary.BigEndian.Uint32(head[0:])
	channel = binary.BigEndian.Uint32(head[4:])
	if length > maxPayload {
		return 0, nil, fmt.Errorf("antui: the server sent %d bytes at once", length)
	}
	payload = make([]byte, length)
	if _, err := io.ReadFull(p.conn, payload); err != nil {
		return 0, nil, err
	}
	return channel, payload, nil
}
