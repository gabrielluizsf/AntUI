package antui

// Touches is every finger on the screen this frame, in the order they
// landed. The slice is the window's own and is rewritten every frame, so
// anything kept has to be copied.
func (win *Window) Touches() []Touch { return win.touches }
