package antui

// Resizable reports whether the window may be resized. It says what was
// asked for, not what the window manager will do about it.
func (win *Window) Resizable() bool { return win != nil && !win.limits.Fixed }
