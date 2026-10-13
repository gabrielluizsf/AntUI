//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"github.com/gabrielluizsf/antui/backend"
)

// handleConfigure records where the window is and, when it is a real resize,
// hands the new size to the canvas.
func (x *Driver) handleConfigure(packet []byte) {
	// Where the window is, which a drag needs: XDND gives the pointer's
	// place on the root and everything here counts from the window.
	x.windowX = int(int16(binary.LittleEndian.Uint16(packet[16:])))
	x.windowY = int(int16(binary.LittleEndian.Uint16(packet[18:])))

	w := int(binary.LittleEndian.Uint16(packet[20:]))
	h := int(binary.LittleEndian.Uint16(packet[22:]))
	cv := x.win.Canvas()
	if w > 0 && h > 0 && (w != cv.Width || h != cv.Height) {
		if x.win.ResizeCanvas(w, h) {
			x.win.Push(backend.Event{Type: backend.EventResize, Width: w, Height: h})
		}
	}
}
