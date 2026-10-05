package antui

// DisplayScale is how many pixels one point of this window covers — the raw
// answer the display gave, and false when it gave none. [Window.ContentScale]
// is the same number with 1 stood in for "did not say", which is what
// arithmetic on a scale wants; the resolution media feature wants the
// difference, because a system that never spoke is one with no resolution to
// report.
//
// The scale is what the system said when the window opened. It is not asked
// again every frame: on X11 that answer costs a round trip to the server.
func (win *Window) DisplayScale() (float64, bool) {
	if win == nil || win.scale <= 0 {
		return 0, false
	}
	return win.scale, true
}

// SystemDark reports the color scheme the system paints its own interface in,
// and whether the system had an answer to give. It follows the OS while the
// window runs — a theme change reaches the next frame — and false is a
// platform with no standard way to ask. prefers-color-scheme then reads the
// condition as unanswered, so the rule stands instead of being guessed at a
// theme the canvas never learned.
//
// [Window.SetSystemDark] stands in for the system.
func (win *Window) SystemDark() (dark, ok bool) {
	if win == nil {
		return false, false
	}
	return win.dark, win.darkKnown
}

// SetSystemDark fixes what [Window.SystemDark] answers, for a platform that
// cannot ask — an X11 window has no standard signal for it — and for a test
// that needs a dark system to draw against. What it is given wins over the
// system for the rest of the window's life, and the next frame draws with it.
func (win *Window) SetSystemDark(dark bool) {
	if win == nil {
		return
	}
	win.dark, win.darkKnown, win.darkForced = dark, true, true
}

// SetDisplayScale fixes what [Window.DisplayScale] answers, for a test that
// needs a dense display to draw against. Zero takes it back to "the display
// did not say".
//
// It is a statement about the display, not a resize: the window keeps the
// size it has, and what changes is the number the resolution media feature
// reads — along with [Window.ContentScale], which is the same number with a
// one stood in for an answer the display never gave.
func (win *Window) SetDisplayScale(scale float64) {
	if win == nil {
		return
	}
	win.scale = scale
}
