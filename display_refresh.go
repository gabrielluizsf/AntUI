package antui

// DisplayRefresh is how many times a second the display refreshes, or 0 when
// it could not be told.
func (win *Window) DisplayRefresh() int {
	if win.native == nil {
		return 0
	}
	return win.native.displayRefresh()
}
