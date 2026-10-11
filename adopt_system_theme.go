package antui

// adoptSystemTheme gives the window the theme the system paints in, when the
// platform had an answer, so that a program which never chooses one is not a
// light window on a dark desktop. It runs once, at open: a window already
// has the light theme from OpenWith, so only a dark system is worth writing
// over it.
//
// From there the theme is the program's to hold — SetTheme and the live
// fields Theme hands out go on from this — while prefers-color-scheme keeps
// following the system frame by frame, on its own.
func (win *Window) adoptSystemTheme() {
	if win == nil || win.native == nil {
		return
	}
	if dark, ok := win.native.systemDark(); ok && dark {
		win.theme = DarkTheme()
	}
}
