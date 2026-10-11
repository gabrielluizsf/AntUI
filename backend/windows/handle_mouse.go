//go:build windows

package windows

// handleMouse folds a pointer message into the window.
func (d *Driver) handleMouse(m wndMsg) (uintptr, bool) {
	switch m.message {
	case wmMouseMove:
		return d.mouseMove(m), true

	case wmMouseLeave:
		d.trackingMouse = false
		return 0, true

	case wmLButtonDown, wmLButtonUp, wmRButtonDown, wmRButtonUp,
		wmMButtonDown, wmMButtonUp:
		return d.mouseButton(m), true

	case wmMouseWheel:
		return d.mouseWheel(m.wparam), true
	}
	return 0, false
}
