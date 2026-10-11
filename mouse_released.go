package antui

// MouseReleased reports whether a button came up during this frame.
func (win *Window) MouseReleased(b MouseButton) bool {
	return b >= 0 && b < mouseCount && win.mouseUpFrame[b]
}
