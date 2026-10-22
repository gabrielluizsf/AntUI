//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// putBitmap sends a whole depth-1 image to drawable through gc in one
// PutImage, format Bitmap — a cursor's source or mask is small enough that
// the chunking a frame's own pixels need never applies here.
func (x *Driver) putBitmap(drawable, gc uint32, bits []byte, width, height int) bool {
	unpadded := 24 + len(bits)
	total := unpadded + pad(unpadded)
	request := make([]byte, total)
	request[0] = xPutImage
	request[1] = 0 // Bitmap
	binary.LittleEndian.PutUint16(request[2:], uint16(total/4))
	binary.LittleEndian.PutUint32(request[4:], drawable)
	binary.LittleEndian.PutUint32(request[8:], gc)
	binary.LittleEndian.PutUint16(request[12:], uint16(width))
	binary.LittleEndian.PutUint16(request[14:], uint16(height))
	binary.LittleEndian.PutUint16(request[16:], 0) // dst-x
	binary.LittleEndian.PutUint16(request[18:], 0) // dst-y
	request[20] = 0                                // left-pad
	request[21] = 1                                // depth
	copy(request[24:], bits)
	return x.write(request)
}
