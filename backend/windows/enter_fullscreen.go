//go:build windows

package windows

import "unsafe"

// enterFullscreen takes the border off and covers the monitor the window is
// mostly on.
func (d *Driver) enterFullscreen(style uintptr) bool {
	// Remember the size to come back out to. The core has recorded it in
	// the canvas, which is about to change; this is the one time it is
	// still the windowed size.
	if cv := d.win.Canvas(); cv != nil && cv.Width > 0 {
		d.windowedWidth, d.windowedHeight = cv.Width, cv.Height
	}

	monitor, _, _ := procMonitorFromWindow.Call(d.hwnd, monitorNearest)
	info := monitorInfo{Size: uint32(unsafe.Sizeof(monitorInfo{}))}
	if ok, _, _ := procGetMonitorInfoW.Call(monitor, uintptr(unsafe.Pointer(&info))); ok == 0 {
		return false
	}
	procSetWindowLongPtr.Call(d.hwnd, gwlStyle,
		style&^uintptr(wsOverlappedWindow)|wsPopup)
	procSetWindowPos.Call(d.hwnd, hwndTop,
		uintptr(info.RcMonitor.Left), uintptr(info.RcMonitor.Top),
		uintptr(info.RcMonitor.Right-info.RcMonitor.Left),
		uintptr(info.RcMonitor.Bottom-info.RcMonitor.Top),
		swpNoOwnerZOrder|swpFrameChanged)
	d.fullscreen = true
	return true
}
