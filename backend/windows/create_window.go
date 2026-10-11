//go:build windows

package windows

import (
	"fmt"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend"
)

// createWindow makes the actual window. The size asked for is the drawable
// area, so the frame is added on top, and the class and handle are the ones
// Open has already prepared.
func createWindow(d *Driver, opts backend.Options) (uintptr, error) {
	r := rect{0, 0, int32(opts.Width), int32(opts.Height)}
	procAdjustWindowRect.Call(uintptr(unsafe.Pointer(&r)), wsOverlappedWindow, 0)

	hwnd, _, err := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr(className))),
		uintptr(unsafe.Pointer(utf16Ptr(opts.Title))),
		wsOverlappedWindow,
		cwUseDefault, cwUseDefault,
		uintptr(r.Right-r.Left), uintptr(r.Bottom-r.Top),
		0, 0, d.instance, 0,
	)
	if hwnd == 0 {
		return 0, fmt.Errorf("antui: could not create the window: %w", err)
	}
	return hwnd, nil
}
