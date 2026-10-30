//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// The ChangeWindowAttributes value-mask bits for a window's background: a
// pixmap — None being zero, which leaves exposed pixels as they are — or a
// solid pixel the server fills them with.
const (
	cwBackPixmap = 0x1
	cwBackPixel  = 0x2
)

// backgroundPixel is the grey the server fills the window with before the
// first frame lands.
const backgroundPixel = 0xFF1F1F1F

// SetExternal says whether another renderer owns the window. While it does,
// the window has no background, so an expose or a resize does not have the
// server fill grey over that renderer's frame; when it stops, the grey comes
// back. It reports whether the request could be made.
func (x *Driver) SetExternal(on bool) bool {
	if !x.alive {
		return false
	}
	x.external = on
	body := make([]byte, 12)
	binary.LittleEndian.PutUint32(body[0:], x.window)
	if on {
		binary.LittleEndian.PutUint32(body[4:], cwBackPixmap)
		binary.LittleEndian.PutUint32(body[8:], 0) // None
	} else {
		binary.LittleEndian.PutUint32(body[4:], cwBackPixel)
		binary.LittleEndian.PutUint32(body[8:], backgroundPixel)
	}
	x.send(xChangeWindowAttributes, 0, body)
	return x.alive
}
