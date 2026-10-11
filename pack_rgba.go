//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// packRGBA swizzles a canvas color into the 0xAABBGGRR word Android's
// RGBA_8888 format expects on a little-endian machine.
func packRGBA(c canvas.Color) uint32 {
	v := uint32(c)
	return v&0xFF00FF00 | (v&0x00FF0000)>>16 | (v&0x000000FF)<<16
}
