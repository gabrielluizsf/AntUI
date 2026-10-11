//go:build windows

package windows

import "unsafe"

// handleResizeMessage answers the two sizing messages that hand the window
// procedure a pointer to write through, which is the whole of the reply.
func (d *Driver) handleResizeMessage(m wndMsg) (uintptr, bool) {
	switch m.message {
	case wmGetMinMaxInfo:
		// The lparam of this message is a pointer the window manager owns and
		// expects to be written through — that is the whole message. Turning
		// it back into one is the only way to answer, so `go vet` flags this
		// line as a possible misuse of unsafe.Pointer and is wrong to: the
		// value came from the OS as a pointer and is live for this call.
		d.minMax((*minMaxInfo)(unsafe.Pointer(m.lparam))) //lint:ignore unsafeptr the OS handed us this pointer
		return 0, true

	case wmSizing:
		d.sizing(m.wparam, (*rect)(unsafe.Pointer(m.lparam))) //lint:ignore unsafeptr the OS handed us this pointer
		return 1, true
	}
	return 0, false
}
