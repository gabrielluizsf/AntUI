package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// BlitScaled draws another canvas fitted into a rectangle.
func (win *Window) BlitScaled(x, y, width, height int, src *canvas.Canvas) {
	win.cv.BlitScaled(x, y, width, height, src)
}
