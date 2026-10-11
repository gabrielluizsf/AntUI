package font

import "encoding/binary"

func be32(b []byte, i int) uint32 {
	if i+4 > len(b) {
		return 0
	}
	return binary.BigEndian.Uint32(b[i:])
}
