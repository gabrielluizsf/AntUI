//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"

	"github.com/gabrielluizsf/antui/backend"
)

// handleDrop answers one of the drag messages. It reports whether the message
// was one of them.
func (x *Driver) handleDrop(packet []byte) bool {
	x.dropAtoms()
	kind := binary.LittleEndian.Uint32(packet[8:])
	data := packet[12:]

	switch kind {
	case x.atomXdndEnter:
		x.dragFrom = binary.LittleEndian.Uint32(data[0:])
		return true

	case x.atomXdndPosition:
		// Where the pointer is, so that a window with more than one place to
		// drop into knows which it is over. The coordinates are on the root,
		// so they are turned into window ones the same way the pointer's are.
		at := binary.LittleEndian.Uint32(data[8:])
		x.mouseX = int(int16(at>>16)) - x.windowX
		x.mouseY = int(int16(at&0xFFFF)) - x.windowY
		x.win.SetMouse(x.mouseX, x.mouseY)
		x.win.Push(backend.Event{Type: backend.EventMouseMove, X: x.mouseX, Y: x.mouseY})

		x.dragFrom = binary.LittleEndian.Uint32(data[0:])
		x.sendDragStatus(true)
		return true

	case x.atomXdndLeave:
		x.dragFrom = 0
		return true

	case x.atomXdndDrop:
		x.dragFrom = binary.LittleEndian.Uint32(data[0:])
		when := binary.LittleEndian.Uint32(data[8:])
		files := x.askSelection(x.atomXdndSelection, x.atomUriList, when)
		x.finishDrag()
		if paths := parseURIList(files); len(paths) > 0 {
			x.win.Push(backend.Event{Type: backend.EventDropFiles, Files: paths,
				X: x.mouseX, Y: x.mouseY})
		}
		return true
	}
	return false
}
