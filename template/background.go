package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) Background(win *antui.Window) canvas.Color {
	return win.Theme().Background
}
