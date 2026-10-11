package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Text draws a string and returns the width of the widest line drawn.
func (win *Window) Text(x, y int, text string, c canvas.Color) int {
	return win.cv.Text(x, y, text, c)
}
