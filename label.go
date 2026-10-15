package antui

import "github.com/gabrielluizsf/antui/canvas"

// Label draws plain text in the theme's text colour.
func (win *Window) Label(x, y int, text string) {
	win.cv.Text(canvas.Point{X: float64(x), Y: float64(y)}, text, win.theme.Text)
}
