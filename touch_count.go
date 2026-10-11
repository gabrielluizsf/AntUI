package antui

// TouchCount is how many fingers are on the screen.
func (win *Window) TouchCount() int { return len(win.touches) }
