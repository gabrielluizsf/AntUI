//go:build windows

package windows

import (
	"github.com/gabrielluizsf/antui/backend"
	"unsafe"
)

// Pump handles whatever happened this frame, folding it into the window
// through win.Push.
func (d *Driver) Pump(win backend.Face) {
	var m msg
	for {
		got, _, _ := procPeekMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0, pmRemove)
		if got == 0 {
			return
		}
		if m.Message == wmQuit {
			win.SetShouldClose()
			return
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}
