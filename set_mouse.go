package antui

// SetMouse is where a backend with no cursor position of its own records the
// pointer.
func (win *Window) SetMouse(x, y int) { win.mouseX, win.mouseY = x, y }
