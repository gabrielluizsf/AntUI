package antui

// SetFixedSize locks the window to one size: both the smallest and the
// largest it may be, with the maximise button gone.
func (win *Window) SetFixedSize(width, height int) {
	if win == nil || width < 1 || height < 1 {
		return
	}
	limits := win.limits
	limits.Fixed, limits.NoMaximize = true, true
	limits.MinWidth, limits.MaxWidth = width, width
	limits.MinHeight, limits.MaxHeight = height, height
	win.SetLimits(limits)
}
