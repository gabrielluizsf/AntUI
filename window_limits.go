package antui

// Limits returns the constraints the window was given.
func (win *Window) Limits() Limits {
	if win == nil {
		return Limits{}
	}
	return win.limits
}
