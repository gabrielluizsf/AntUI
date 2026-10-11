//go:build windows

package windows

import (
	"unsafe"
)

// SetTitle names the window in the system's decoration.
func (d *Driver) SetTitle(title string) {
	if d.hwnd != 0 {
		procSetWindowTextW.Call(d.hwnd, uintptr(unsafe.Pointer(utf16Ptr(title))))
	}
}
