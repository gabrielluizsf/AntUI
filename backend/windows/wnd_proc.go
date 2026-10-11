//go:build windows

package windows

func wndProc(hwnd, message, wparam, lparam uintptr) uintptr {
	windowsMu.RLock()
	d := windows[hwnd]
	windowsMu.RUnlock()

	// Messages arrive before CreateWindowExW has returned the handle we key
	// the map on, so there is a window of time with nothing to dispatch to.
	if d == nil || d.win == nil {
		return defWindowProc(hwnd, message, wparam, lparam)
	}
	if ret, handled := d.handleMessage(wndMsg{hwnd, message, wparam, lparam}); handled {
		return ret
	}
	return defWindowProc(hwnd, message, wparam, lparam)
}
