package antui

// MouseDown reports whether a button is held.
func (win *Window) MouseDown(b MouseButton) bool {
	return b >= 0 && b < mouseCount && win.mouseState[b]
}
