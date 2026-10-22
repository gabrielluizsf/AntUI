//go:build windows

package windows

// htClient is the hit-test code WM_SETCURSOR carries in the low word of its
// lParam when the pointer is over the window's own content, rather than a
// border, a resize handle, or some other part Windows draws itself.
const htClient = 1

// handleSetCursor keeps the custom cursor on the pointer: Windows sends this
// message every time the mouse moves over the window, and forgets the
// picture that was there a moment ago unless something sets it again.
func (d *Driver) handleSetCursor(m wndMsg) (uintptr, bool) {
	if loWord(m.lparam) != htClient {
		return 0, false
	}
	procSetCursor.Call(d.cursorOrArrow())
	return 1, true
}
