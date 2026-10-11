package antui

// Background paints the window in the theme's background colour.
func (win *Window) Background() { win.cv.Clear(win.theme.Background) }
