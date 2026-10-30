//go:build windows

package windows

import (
	"github.com/gabrielluizsf/antui/backend"
	"unsafe"
)

// paint answers WM_PAINT: the freshly drawn frame is blitted, the paint is
// closed, and the window is told to redraw. While another renderer owns the
// window the blit is skipped — the paint still has to be begun and ended, or
// Windows sends WM_PAINT again forever.
func (d *Driver) paint(hwnd uintptr) {
	var ps paintStruct
	hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
	if !d.external {
		d.blit(hdc)
	}
	procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
	d.win.Push(backend.Event{Type: backend.EventExpose})
}
