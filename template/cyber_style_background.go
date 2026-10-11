package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (cyberStyle) Background(win *antui.Window) canvas.Color {
	return cyberPalette.bg
}
