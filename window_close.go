package antui

// Close destroys the window and gives back everything it held.
func (win *Window) Close() {
	if win == nil || win.native == nil {
		return
	}
	win.native.close()
	win.native = nil
}
