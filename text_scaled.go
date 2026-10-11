package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// TextScaled draws a string with the font blown up, returning the width of
// the widest line drawn.
func (win *Window) TextScaled(x, y int, text string, c canvas.Color, scale int) int {
	return win.cv.TextScaled(x, y, text, c, scale)
}
