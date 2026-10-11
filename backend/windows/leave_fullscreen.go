//go:build windows

package windows

import "unsafe"

// leaveFullscreen puts the frame and the windowed placement back.
func (d *Driver) leaveFullscreen(style uintptr) bool {
	width, height := d.windowedWidth, d.windowedHeight
	if width <= 0 {
		width = 640
	}
	if height <= 0 {
		height = 480
	}
	procSetWindowLongPtr.Call(d.hwnd, gwlStyle,
		style&^uintptr(wsPopup)|wsOverlappedWindow)
	r := rect{0, 0, int32(width), int32(height)}
	procAdjustWindowRect.Call(uintptr(unsafe.Pointer(&r)), wsOverlappedWindow, 0)
	procSetWindowPos.Call(d.hwnd, hwndNoTopMost, cwUseDefault, cwUseDefault,
		uintptr(r.Right-r.Left), uintptr(r.Bottom-r.Top),
		swpNoMove|swpNoOwnerZOrder|swpFrameChanged)
	d.fullscreen = false
	return true
}
