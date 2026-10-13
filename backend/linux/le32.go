//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

func le32(values ...uint32) []byte {
	buf := make([]byte, 4*len(values))
	for i, v := range values {
		binary.LittleEndian.PutUint32(buf[i*4:], v)
	}
	return buf
}
