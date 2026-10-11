package antui

// Quit marks the window for closing, so the next Begin returns false.
func (win *Window) Quit() { win.shouldClose = true }
