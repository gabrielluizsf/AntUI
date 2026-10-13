//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// request builds a request whose length field it fills in itself. The body
// is padded to a four-byte boundary, as every request must be.
func (x *Driver) request(opcode, detail byte, body []byte) []byte {
	total := 4 + len(body)
	total += pad(total)
	buf := make([]byte, total)
	buf[0] = opcode
	buf[1] = detail
	binary.LittleEndian.PutUint16(buf[2:], uint16(total/4))
	copy(buf[4:], body)
	return buf
}
