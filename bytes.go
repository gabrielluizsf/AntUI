//go:build linux && !android

package antui

import "encoding/binary"

func (t *tags) bytes(b []byte) {
	t.b = append(t.b, tagArbitrary)
	t.b = binary.BigEndian.AppendUint32(t.b, uint32(len(b)))
	t.b = append(t.b, b...)
}
