package antui

// SetFullscreen asks the system to put the window full screen, or to take it
// back out. It reports whether the request could be made, not whether it was
// granted.
func (win *Window) SetFullscreen(on bool) bool {
	if on && !win.fullscreen {
		win.windowedWidth = win.cv.Width
		win.windowedHeight = win.cv.Height
	}
	if win.native == nil || !win.native.setFullscreen(on) {
		return false
	}
	win.fullscreen = on
	win.fullRedraw = true
	return true
}
