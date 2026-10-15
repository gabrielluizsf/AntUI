package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// TextScaled draws a string with the font blown up, returning the width of
// the widest line drawn.
func (win *Window) TextScaled(p canvas.Point, params TextScaledParams) int {
	return win.cv.TextScaled(p, params.Text, canvas.TextStyle{Color: params.Color, Scale: params.Scale})
}
