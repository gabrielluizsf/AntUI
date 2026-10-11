package antui

// bounds resolves the limits against the size the window is now: a fixed
// window's smallest and largest size are both the size it has. Keeping Fixed
// as a flag rather than writing the numbers in is what lets a window locked
// at no particular size follow a display change instead of being frozen at
// whatever it happened to be when the game asked.
func (win *Window) Bounds() Limits {
	l := win.limits
	if l.Fixed {
		// A lock that names no size locks the window at whatever it is now.
		// One that names a size is that size, and the maximum is the half
		// believed when the two disagree — the same rule clamp follows.
		width := orElse(l.MaxWidth, orElse(l.MinWidth, win.width))
		height := orElse(l.MaxHeight, orElse(l.MinHeight, win.height))
		l.MinWidth, l.MaxWidth = width, width
		l.MinHeight, l.MaxHeight = height, height
	}
	return l
}
