//go:build windows

package windows

import "github.com/gabrielluizsf/antui/backend"

// resize follows a WM_SIZE: the canvas grows to the new drawable area and the
// window is told, unless the size is the one already held.
func (d *Driver) resize(win backend.Face, wparam, lparam uintptr) {
	width, height := int(uint16(lparam)), int(uint16(lparam>>16))
	cv := win.Canvas()
	if wparam != sizeMinimized && width > 0 && height > 0 &&
		(width != cv.Width || height != cv.Height) {
		if win.ResizeCanvas(width, height) {
			win.Push(backend.Event{Type: backend.EventResize, Width: width, Height: height})
		}
	}
}
