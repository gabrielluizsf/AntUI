package antui

// SetLimits changes what the window may be resized to, and takes effect at
// once: a window already outside the new bounds is asked to come back inside
// them.
func (win *Window) SetLimits(limits Limits) {
	if win == nil {
		return
	}
	win.limits = limits
	if win.native != nil {
		win.native.setLimits(win, win.Bounds())
	}
	if width, height := limits.Clamp(win.width, win.height); width != win.width ||
		height != win.height {
		win.SetSize(width, height)
	}
}
