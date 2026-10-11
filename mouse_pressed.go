package antui

// MousePressed reports whether a button went down during this frame.
func (win *Window) MousePressed(b MouseButton) bool {
	return b >= 0 && b < mouseCount && win.mouseDownFrame[b]
}
