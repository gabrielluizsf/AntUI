package antui

// Label draws plain text in the theme's text colour.
func (win *Window) Label(x, y int, text string) {
	win.cv.Text(x, y, text, win.theme.Text)
}
