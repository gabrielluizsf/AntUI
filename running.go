package antui

// Running reports whether the window is still open.
func (win *Window) Running() bool { return win != nil && !win.shouldClose }
