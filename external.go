package antui

// External reports whether another renderer owns the window's surface; see
// [Window.SetExternal].
func (win *Window) External() bool { return win.external }
