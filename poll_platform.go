package antui

// pollPlatform pumps the backend and refreshes the system theme and display
// scale.
//
// The system's color scheme and the display's scale are both asked every
// frame rather than remembered at open: a theme or a scale the platform told
// us about is what prefers-color-scheme and the resolution feature read on
// the frame that follows it. Neither costs a round trip — the platform keeps
// its own answer and hands it over — and an override set by SetSystemDark or
// SetDisplayScale stands over both.
func (win *Window) pollPlatform() {
	if win.native == nil {
		return
	}
	win.native.pump(win)
	if !win.darkForced {
		win.dark, win.darkKnown = win.native.systemDark()
	}
	if !win.scaleForced {
		// An unanswered read leaves what the display said before: a scale is
		// a property of the screen, and a read that got no answer is not the
		// display taking the number back.
		if scale, ok := win.native.displayScale(); ok {
			win.scale = scale
		}
	}
}
