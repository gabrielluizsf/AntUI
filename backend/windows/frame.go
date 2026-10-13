//go:build windows

package windows

import "unsafe"

// frame is how much wider and taller the whole window is than its drawable
// area, which every one of these calls is in terms of and none of them says.
func (d *Driver) frame() (dx, dy int32) {
	style, _, _ := procGetWindowLongPtr.Call(d.hwnd, gwlStyle)
	r := rect{0, 0, 100, 100}
	procAdjustWindowRect.Call(uintptr(unsafe.Pointer(&r)), style, 0)
	return (r.Right - r.Left) - 100, (r.Bottom - r.Top) - 100
}
