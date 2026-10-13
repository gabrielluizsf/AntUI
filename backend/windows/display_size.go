//go:build windows

package windows

import "unsafe"

// DisplaySize is the size of the display the window is on.
func (d *Driver) DisplaySize() (w, h int, ok bool) {
	if d.hwnd != 0 {
		monitor, _, _ := procMonitorFromWindow.Call(d.hwnd, monitorNearest)
		info := monitorInfo{Size: uint32(unsafe.Sizeof(monitorInfo{}))}
		if ok, _, _ := procGetMonitorInfoW.Call(monitor, uintptr(unsafe.Pointer(&info))); ok != 0 {
			return int(info.RcMonitor.Right - info.RcMonitor.Left),
				int(info.RcMonitor.Bottom - info.RcMonitor.Top), true
		}
	}
	width, _, _ := procGetSystemMetrics.Call(smCXScreen)
	height, _, _ := procGetSystemMetrics.Call(smCYScreen)
	if width == 0 || height == 0 {
		return 0, 0, false
	}
	return int(width), int(height), true
}
