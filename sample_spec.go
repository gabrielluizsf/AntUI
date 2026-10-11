//go:build linux && !android

package antui

import (
	"encoding/binary"
)

func (t *tags) sampleSpec(format, channels byte, rate uint32) {
	t.b = append(t.b, tagSampleSpec, format, channels)
	t.b = binary.BigEndian.AppendUint32(t.b, rate)
}
