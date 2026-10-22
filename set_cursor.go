package antui

// SetCursor replaces the mouse pointer with a custom picture, reporting
// whether the window system took it. A zero-value icon — Image nil — puts
// the system's own pointer back.
func (win *Window) SetCursor(icon CursorIcon) bool {
	if win == nil || win.native == nil {
		return false
	}
	return win.native.setCursor(icon)
}
