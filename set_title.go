package antui

// SetTitle changes the title bar.
func (win *Window) SetTitle(title string) {
	if win.native != nil {
		win.native.setTitle(title)
	}
}
