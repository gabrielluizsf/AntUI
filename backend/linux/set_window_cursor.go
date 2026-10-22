//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// cwCursor is the ChangeWindowAttributes value-mask bit that says a cursor
// value follows.
const cwCursor = 0x4000

// setWindowCursor makes cursor the window's own, or — cursor zero, the None
// value every core request reads it as — gives up the window's cursor and
// falls back to its parent's.
func (x *Driver) setWindowCursor(cursor uint32) {
	body := make([]byte, 12)
	binary.LittleEndian.PutUint32(body[0:], x.window)
	binary.LittleEndian.PutUint32(body[4:], cwCursor)
	binary.LittleEndian.PutUint32(body[8:], cursor)
	x.send(xChangeWindowAttributes, 0, body)
}
