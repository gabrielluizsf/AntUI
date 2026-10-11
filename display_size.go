package antui

// DisplaySize is the display's size in pixels, and whether it could be told.
func (win *Window) DisplaySize() (width, height int, ok bool) {
	if win.native == nil {
		return 0, 0, false
	}
	return win.native.displaySize()
}
