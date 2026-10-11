//go:build windows

package windows

import (
	"github.com/gabrielluizsf/antui/backend"
	"unsafe"
)

// mouseMove records the pointer's new place and asks to be told when it
// leaves, so the tracking flag can be reset and the request made again next
// time it comes back.
func (d *Driver) mouseMove(m wndMsg) uintptr {
	d.mouseX, d.mouseY = loWord(m.lparam), hiWord(m.lparam)
	if !d.trackingMouse {
		track := trackMouseEvent{
			Size:      uint32(unsafe.Sizeof(trackMouseEvent{})),
			Flags:     tmeLeave,
			HWndTrack: m.hwnd,
		}
		procTrackMouseEvent.Call(uintptr(unsafe.Pointer(&track)))
		d.trackingMouse = true
	}
	d.win.Push(backend.Event{
		Type: backend.EventMouseMove,
		X:    d.mouseX,
		Y:    d.mouseY,
		Mods: winMods(),
	})
	return 0
}
