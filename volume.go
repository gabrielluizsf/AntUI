//go:build linux && !android

package antui

import (
	"encoding/binary"
)

func (t *tags) volume(channels int, level uint32) {
	t.b = append(t.b, tagCvolume, byte(channels))
	for range channels {
		t.b = binary.BigEndian.AppendUint32(t.b, level)
	}
}
