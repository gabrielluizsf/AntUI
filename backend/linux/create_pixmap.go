//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// createPixmap makes a pixmap of the given depth, the size of width by
// height, and returns its id.
func (x *Driver) createPixmap(depth byte, width, height int) uint32 {
	pid := x.newID()
	body := make([]byte, 12)
	binary.LittleEndian.PutUint32(body[0:], pid)
	binary.LittleEndian.PutUint32(body[4:], x.window)
	binary.LittleEndian.PutUint16(body[8:], uint16(width))
	binary.LittleEndian.PutUint16(body[10:], uint16(height))
	x.send(xCreatePixmap, depth, body)
	return pid
}
