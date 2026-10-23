//go:build windows

package windows

import "unsafe"

// removeNotifyIcon takes the window's notification-area icon away, if
// [Driver.SendNotification] ever put one there. Windows does not do it on
// its own: an icon outlives its window until the mouse happens to pass over
// it.
func (d *Driver) removeNotifyIcon() {
	if !d.notifyIcon {
		return
	}
	data := d.notifyIconData()
	procShellNotifyIconW.Call(nimDelete, uintptr(unsafe.Pointer(&data)))
	d.notifyIcon = false
}
