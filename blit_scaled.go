package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// BlitScaled draws another canvas fitted into a rectangle.
func (win *Window) BlitScaled(area canvas.Area, src *canvas.Canvas) {
	win.cv.BlitScaled(area, src)
}
