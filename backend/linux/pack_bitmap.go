//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// packBitmap turns a width-by-height grid of 1s and 0s into the byte layout
// PutImage wants for a depth-1 image: each row padded out to a whole
// [Driver.bitmapScanlinePad]-bit unit, the pixels within a unit ordered by
// [Driver.bitmapBitOrder], and the server's own [Driver.imageByteOrder]
// deciding which byte of a multi-byte unit comes first on the wire — the
// same three answers a frame's own pixels are packed by, asked of a single
// bit instead of thirty-two.
func (x *Driver) packBitmap(on func(px, py int) bool, width, height int) []byte {
	unitBits := int(x.bitmapScanlinePad)
	if unitBits != 8 && unitBits != 16 && unitBits != 32 {
		unitBits = 32
	}
	unitBytes := unitBits / 8
	rowUnits := (width + unitBits - 1) / unitBits
	rowBytes := rowUnits * unitBytes
	out := make([]byte, rowBytes*height)

	msb := x.bitmapBitOrder != 0
	for py := 0; py < height; py++ {
		row := out[py*rowBytes : (py+1)*rowBytes]
		for ui := 0; ui < rowUnits; ui++ {
			var unit uint32
			for bit := 0; bit < unitBits; bit++ {
				px := ui*unitBits + bit
				if px >= width || !on(px, py) {
					continue
				}
				shift := bit
				if msb {
					shift = unitBits - 1 - bit
				}
				unit |= 1 << uint(shift)
			}
			dst := row[ui*unitBytes : (ui+1)*unitBytes]
			switch unitBytes {
			case 1:
				dst[0] = byte(unit)
			case 2:
				if x.imageByteOrder == 0 {
					binary.LittleEndian.PutUint16(dst, uint16(unit))
				} else {
					binary.BigEndian.PutUint16(dst, uint16(unit))
				}
			case 4:
				if x.imageByteOrder == 0 {
					binary.LittleEndian.PutUint32(dst, unit)
				} else {
					binary.BigEndian.PutUint32(dst, unit)
				}
			}
		}
	}
	return out
}
