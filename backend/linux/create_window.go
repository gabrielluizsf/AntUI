//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// createWindow makes the window itself. The depth, visual and colormap are
// all inherited from the parent, which is the root: this window wants exactly
// what the screen already is.
func (x *Driver) createWindow(title string, width, height int) {
	x.window = x.newID()
	body := make([]byte, 36)
	binary.LittleEndian.PutUint32(body[0:], x.window)
	binary.LittleEndian.PutUint32(body[4:], x.root)
	binary.LittleEndian.PutUint16(body[8:], 0)  // x
	binary.LittleEndian.PutUint16(body[10:], 0) // y
	binary.LittleEndian.PutUint16(body[12:], uint16(width))
	binary.LittleEndian.PutUint16(body[14:], uint16(height))
	binary.LittleEndian.PutUint16(body[16:], 0) // border width
	binary.LittleEndian.PutUint16(body[18:], 1) // class = InputOutput
	binary.LittleEndian.PutUint32(body[20:], 0) // visual = the parent's
	binary.LittleEndian.PutUint32(body[24:], 0x2|0x800)
	binary.LittleEndian.PutUint32(body[28:], 0xFF1F1F1F) // background pixel
	binary.LittleEndian.PutUint32(body[32:], xEventMask)
	x.send(xCreateWindow, 0, body) // detail 0 = inherit the parent's depth

	x.setTitleOn(title)
	x.changeProperty(atomWMClass, atomString, 8, []byte("antui\x00AntUI\x00"), 12)

	// WM_HINTS with the input bit: without it some window managers never hand
	// keyboard focus to the window at all.
	hints := le32(1|2, 1, 1, 0, 0, 0, 0, 0, 0) // InputHint|StateHint, input, NormalState
	x.changeProperty(atomWMHints, atomWMHints, 32, hints, 9)

	if x.atomWMProtocols != 0 && x.atomWMDeleteWindow != 0 {
		// Without this the window manager kills the connection to close the
		// window instead of asking, and the program never gets to save.
		x.changeProperty(x.atomWMProtocols, atomAtom, 32, le32(x.atomWMDeleteWindow), 1)
	}
}
