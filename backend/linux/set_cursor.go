//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "github.com/gabrielluizsf/antui/backend"

// SetCursor replaces the mouse pointer with icon's picture, reporting
// whether the request could be made. A nil Image puts the window's own
// cursor back to None, which is the parent's — the window manager's
// ordinary arrow.
func (x *Driver) SetCursor(icon backend.CursorIcon) bool {
	if !x.alive {
		return false
	}
	old := x.cursor
	if icon.Image == nil || icon.Image.Width <= 0 || icon.Image.Height <= 0 {
		x.cursor = 0
	} else {
		x.cursor = x.createCursor(icon)
	}
	x.setWindowCursor(x.cursor)
	if old != 0 {
		x.freeCursor(old)
	}
	return true
}
