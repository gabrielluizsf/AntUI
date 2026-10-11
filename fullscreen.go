package antui

// Fullscreen reports what was last asked for, not what a window manager did
// about it.
func (win *Window) Fullscreen() bool { return win.fullscreen }
