//go:build windows

package windows

// Native is the window's HWND and the HINSTANCE its class was registered
// from, for another renderer to create its surface on. The HWND is zero once
// the window is gone.
func (d *Driver) Native() (hwnd, instance uintptr) { return d.hwnd, d.instance }
