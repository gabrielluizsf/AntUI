//go:build windows

package windows

import (
	"github.com/gabrielluizsf/antui/backend"
)

// SetSize asks the platform for a new drawable size, reporting whether the
// request could be made.
func (d *Driver) SetSize(win backend.Face, width, height int) bool {
	if d.hwnd == 0 {
		return false
	}
	d.windowedWidth, d.windowedHeight = width, height
	dx, dy := d.frame()
	procSetWindowPos.Call(d.hwnd, hwndTop, 0, 0,
		uintptr(int32(width)+dx), uintptr(int32(height)+dy),
		swpNoMove|swpNoZOrder|swpNoActivate)
	return true
}
