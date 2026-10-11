package antui

// SystemDark reports the color scheme the system paints its own interface in,
// and whether the system had an answer to give. It follows the OS while the
// window runs — a theme change reaches the next frame — and false is a
// system with nothing to ask: a display server carries no such signal, and
// where a platform reads one it is reading the desktop's own settings, which
// a desktop may decline to answer. prefers-color-scheme then reads the
// condition as unanswered, so the rule stands instead of being guessed at a
// theme the canvas never learned.
//
// The same answer picks the window's starting theme, once, at open; from
// there [Window.SetTheme] and the fields [Window.Theme] hands out are the
// program's to change.
//
// [Window.SetSystemDark] stands in for the system.
func (win *Window) SystemDark() (dark, ok bool) {
	if win == nil {
		return false, false
	}
	return win.dark, win.darkKnown
}
