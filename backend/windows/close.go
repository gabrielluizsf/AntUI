//go:build windows

package windows

import "runtime"

// Close destroys the window and lets go of everything the platform held.
func (d *Driver) Close() {
	if d.hwnd == 0 {
		return
	}
	windowsMu.Lock()
	delete(windows, d.hwnd)
	windowsMu.Unlock()

	d.removeNotifyIcon()
	procDestroyWindow.Call(d.hwnd)
	d.hwnd = 0
	runtime.UnlockOSThread()
}
