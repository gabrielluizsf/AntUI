package antui

// SetSystemDark fixes what [Window.SystemDark] answers, for a system with
// nothing to say and for a test that needs a dark one to draw against. What
// it is given wins over the system for the rest of the window's life, and
// the next frame draws with it.
func (win *Window) SetSystemDark(dark bool) {
	if win == nil {
		return
	}
	win.dark, win.darkKnown, win.darkForced = dark, true, true
}
