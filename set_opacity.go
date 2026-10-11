package antui

// SetOpacity fades the whole window — decorations and all — so that what sits
// behind the program shows through; 255 is fully opaque and 0 is invisible.
// It reports whether the request could be made: there is no window to fade on
// an offscreen canvas, so that reports false while still remembering the
// value. The window itself is never returned to — this is a whole-window
// opacity, not a fade of the pixels drawn into it, so what is drawn keeps its
// contrast and only its overall brightness changes.
func (win *Window) SetOpacity(alpha uint8) bool {
	win.opacity = alpha
	if win.native == nil {
		return false
	}
	return win.native.setOpacity(alpha)
}
