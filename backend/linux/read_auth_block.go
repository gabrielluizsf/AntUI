//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"io"
)

// readAuthBlock reads one length-prefixed field of an authority file entry.
func readAuthBlock(f io.Reader) ([]byte, error) {
	var length uint16
	if err := binary.Read(f, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	buf := make([]byte, length)
	_, err := io.ReadFull(f, buf)
	return buf, err
}
