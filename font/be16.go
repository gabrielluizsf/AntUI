package font

import "encoding/binary"

func be16(b []byte, i int) uint16 {
	if i+2 > len(b) {
		return 0
	}
	return binary.BigEndian.Uint16(b[i:])
}
