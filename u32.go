//go:build linux && !android

package antui

import "encoding/binary"

func (t *tags) u32(v uint32) {
	t.b = append(t.b, tagU32)
	t.b = binary.BigEndian.AppendUint32(t.b, v)
}
