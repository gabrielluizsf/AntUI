//go:build windows

package windows

import (
	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

// Present copies the window's dirty rectangle to the platform surface.
//
// GDI is not told the rectangle is dirty; it is handed the whole frame the
// way a game hands one over, and the dirty rectangle decides on the other
// side whether the frame is even copied. Which is this call as it is.
func (d *Driver) Present(win backend.Face, dirty canvas.Area) {
	if dirty.Width <= 0 || dirty.Height <= 0 || d.hwnd == 0 {
		return
	}
	hdc, _, _ := procGetDC.Call(d.hwnd)
	if hdc == 0 {
		return
	}
	d.blit(hdc)
	procReleaseDC.Call(d.hwnd, hdc)
}
