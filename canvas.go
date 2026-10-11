package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Canvas is the window's pixel surface: direct access, if you need it.
func (win *Window) Canvas() *canvas.Canvas { return win.cv }
