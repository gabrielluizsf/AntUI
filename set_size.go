package antui

// SetSize asks for a new drawable size, reduced to the limits and to what the
// display can show. It reports whether the request could be made at all, not
// whether the window manager honoured it.
//
// The canvas is not resized here. It is resized when the size actually
// arrives, which is what keeps everything drawn following the framebuffer
// rather than following a request that may have been ignored.
func (win *Window) SetSize(width, height int) bool {
	if win == nil || win.native == nil || width < 1 || height < 1 {
		return false
	}
	if dw, dh, ok := win.native.displaySize(); ok {
		width, height = fitDisplay(width, height, dw, dh)
	}
	width, height = win.limits.Clamp(width, height)
	if width == win.width && height == win.height {
		return true
	}
	return win.native.setSize(win, width, height)
}
