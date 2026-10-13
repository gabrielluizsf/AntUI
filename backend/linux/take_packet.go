//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// takePacket pulls one whole packet out of the buffer: the fixed 32 bytes,
// plus the variable tail a reply carries. It reports false when only part of
// one has arrived, leaving what it has for the next pass.
func (x *Driver) takePacket() (packet, extra []byte, ok bool) {
	if x.available() < xEventPacketBytes {
		return nil, nil, false
	}
	p := x.rbuf[x.rpos:]
	total := xEventPacketBytes
	if p[0] == xReply {
		total += int(binary.LittleEndian.Uint32(p[4:])) * xProtocolWordBytes
	}
	if x.available() < total {
		return nil, nil, false
	}
	// The slices point into the read buffer, which the next fill will move,
	// so both are copied before anything else can touch it.
	packet = append([]byte(nil), x.rbuf[x.rpos:x.rpos+xEventPacketBytes]...)
	if total > xEventPacketBytes {
		extra = append([]byte(nil), x.rbuf[x.rpos+xEventPacketBytes:x.rpos+total]...)
	}
	x.consume(total)
	return packet, extra, true
}
