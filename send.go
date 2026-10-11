//go:build linux && !android

package antui

import (
	"encoding/binary"
)

// send writes one packet: the twenty-byte descriptor, then the payload.
func (p *pulse) send(channel uint32, payload []byte) error {
	var head [20]byte
	binary.BigEndian.PutUint32(head[0:], uint32(len(payload)))
	binary.BigEndian.PutUint32(head[4:], channel)
	binary.BigEndian.PutUint32(head[8:], 0)  // offset, high half
	binary.BigEndian.PutUint32(head[12:], 0) // and low
	binary.BigEndian.PutUint32(head[16:], 0) // flags
	if _, err := p.conn.Write(head[:]); err != nil {
		return err
	}
	_, err := p.conn.Write(payload)
	return err
}
