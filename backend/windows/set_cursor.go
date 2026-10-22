//go:build windows

package windows

import "github.com/gabrielluizsf/antui/backend"

// SetCursor replaces the mouse pointer with icon's picture, reporting
// whether the request could be made. A nil Image puts the window's pointer
// back to the ordinary system arrow.
func (d *Driver) SetCursor(icon backend.CursorIcon) bool {
	if d.hwnd == 0 {
		return false
	}
	old := d.cursor
	if icon.Image == nil || icon.Image.Width <= 0 || icon.Image.Height <= 0 {
		d.cursor = 0
	} else {
		h := hcursor(icon.Image, icon.HotX, icon.HotY)
		if h == 0 {
			return false
		}
		d.cursor = h
	}
	procSetCursor.Call(d.cursorOrArrow())
	if old != 0 {
		procDestroyCursor.Call(old)
	}
	return true
}
