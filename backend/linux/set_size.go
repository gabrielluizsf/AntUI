//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"

	"github.com/gabrielluizsf/antui/backend"
)

// SetSize asks the window manager for a new drawable size, reporting
// whether the request could be made. The window manager may refuse, and the
// canvas follows the ConfigureNotify rather than this.
func (x *Driver) SetSize(win backend.Face, width, height int) bool {
	if !x.alive || x.window == 0 {
		return false
	}
	// window, then a mask of which values follow, then those values: the
	// width and the height and nothing else, so the window manager keeps
	// deciding where the window goes.
	body := make([]byte, 16)
	binary.LittleEndian.PutUint32(body[0:], x.window)
	binary.LittleEndian.PutUint16(body[4:], 0x4|0x8)
	binary.LittleEndian.PutUint32(body[8:], uint32(max(width, 1)))
	binary.LittleEndian.PutUint32(body[12:], uint32(max(height, 1)))
	x.send(xConfigureWindow, 0, body)
	return true
}
