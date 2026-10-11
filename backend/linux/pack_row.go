//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"unsafe"

	"github.com/gabrielluizsf/antui/canvas"
)

// packRow writes one row of pixels in the byte order the server asked for.
//
// The ordinary case is that it asked for the order the pixels are already in
// — a little-endian server on a little-endian machine, with red where red
// goes — and then this is a memmove. It matters more than it looks: a frame
// is a third of a million pixels, and doing them one at a time through a
// function pointer was costing more per frame than drawing the world did.
//
// The slow path is still there and still right, because a server that wants
// the other order is not a server that should get the wrong picture.
func (x *Driver) packRow(dst []byte, src []canvas.Color) {
	if len(src) == 0 {
		return
	}
	if x.samePixels {
		copy(dst, unsafe.Slice((*byte)(unsafe.Pointer(&src[0])), len(src)*4))
		return
	}

	put := binary.LittleEndian.PutUint32
	if x.imageByteOrder != 0 {
		put = binary.BigEndian.PutUint32
	}
	for i, pixel := range src {
		c := uint32(pixel)
		if x.swapRB {
			c = c&0xFF00FF00 | (c&0xFF)<<16 | (c>>16)&0xFF
		}
		put(dst[i*4:], c)
	}
}
