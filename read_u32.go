//go:build linux && !android

package antui

import (
	"encoding/binary"
)

func (r *read) u32() (uint32, error) {
	if err := r.tag(tagU32); err != nil {
		return 0, err
	}
	if r.at+4 > len(r.b) {
		return 0, errShort
	}
	v := binary.BigEndian.Uint32(r.b[r.at:])
	r.at += 4
	return v, nil
}
