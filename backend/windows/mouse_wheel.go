//go:build windows

package windows

import "github.com/gabrielluizsf/antui/backend"

// mouseWheel pushes one wheel step: Windows reports the delta in the high
// word of wparam, and the window counts in steps of its own.
func (d *Driver) mouseWheel(wparam uintptr) uintptr {
	d.win.Push(backend.Event{
		Type:  backend.EventMouseWheel,
		Wheel: int(int16(uint16(wparam>>16))) / wheelDelta,
		X:     d.mouseX,
		Y:     d.mouseY,
		Mods:  winMods(),
	})
	return 0
}
