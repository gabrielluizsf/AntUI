package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Blit composites another canvas over the window.
func (win *Window) Blit(x, y int, src *canvas.Canvas) { win.cv.Blit(x, y, src) }
