//go:build windows

package windows

import (
	"github.com/gabrielluizsf/antui/backend"
	"unsafe"
)

// paint answers WM_PAINT: the freshly drawn frame is blitted, the paint is
// closed, and the window is told to redraw.
func (d *Driver) paint(hwnd uintptr) {
	var ps paintStruct
	hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
	d.blit(hdc)
	procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
	d.win.Push(backend.Event{Type: backend.EventExpose})
}
