package antui

// focusLost clears the held keys and buttons and ends every finger: the
// window is not being typed at or clicked any more.
func (win *Window) focusLost() {
	win.keyState = [keyCount]bool{}
	win.mouseState = [mouseCount]bool{}
	win.mods = 0
	for i := range win.touches {
		win.touches[i].Ended = true
		win.touches[i].Cancelled = true
	}
}
