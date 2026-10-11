//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// createGC makes the graphics context every drawing request borrows: the
// window, the screen default, and nothing else set.
func (x *Driver) createGC() {
	x.gc = x.newID()
	gcBody := make([]byte, 12)
	binary.LittleEndian.PutUint32(gcBody[0:], x.gc)
	binary.LittleEndian.PutUint32(gcBody[4:], x.window)
	x.send(xCreateGC, 0, gcBody)
}
