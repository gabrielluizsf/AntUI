//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"

	"github.com/gabrielluizsf/antui/canvas"
)

// putChunk sends one PutImage request for the rows staged in scratch,
// reporting whether the write got through.
func (x *Driver) putChunk(dirty canvas.Area, y, rows, rowBytes int) bool {
	header := x.header[:20:20]
	binary.LittleEndian.PutUint32(header[0:], x.window)
	binary.LittleEndian.PutUint32(header[4:], x.gc)
	binary.LittleEndian.PutUint16(header[8:], uint16(dirty.Width))
	binary.LittleEndian.PutUint16(header[10:], uint16(rows))
	binary.LittleEndian.PutUint16(header[12:], uint16(int16(dirty.X)))
	binary.LittleEndian.PutUint16(header[14:], uint16(int16(dirty.Y+y)))
	header[16] = 0 // left pad
	header[17] = x.rootDepth

	// PutImage carries the pixels as its body, so the request is built
	// here rather than through send: the data is already padded to four
	// bytes by being whole pixels.
	x.seq++
	total := 24 + rowBytes*rows

	// The request buffer is kept rather than made: at a megabyte and a
	// half a frame, allocating it is a garbage collection every couple of
	// seconds and a pause in the middle of somebody's game.
	if cap(x.putBuffer) < total {
		x.putBuffer = make([]byte, total)
	}
	request := x.putBuffer[:total]
	request[0] = xPutImage
	request[1] = 2 // ZPixmap
	binary.LittleEndian.PutUint16(request[2:], uint16(total/4))
	copy(request[4:], header)
	copy(request[24:], x.scratch[:rowBytes*rows])

	return x.write(request)
}
