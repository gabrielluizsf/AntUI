package antui

// SetDisplayScale fixes what [Window.DisplayScale] answers, for a test that
// needs a dense display to draw against. Zero takes it back to "the display
// did not say". What it is given wins over the display for the rest of the
// window's life, so the frame stops asking it.
//
// It is a statement about the display, not a resize: the window keeps the
// size it has, and what changes is the number the resolution media feature
// reads — along with [Window.ContentScale], which is the same number with a
// one stood in for an answer the display never gave.
func (win *Window) SetDisplayScale(scale float64) {
	if win == nil {
		return
	}
	win.scale, win.scaleForced = scale, true
}
