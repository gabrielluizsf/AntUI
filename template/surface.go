package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// surface is the colour a reversed-out number draws in on a selected day
// chip: the widget's own background when it has one, or the theme's
// background otherwise — which, against the theme's text (the chip's ink),
// always contrasts.
func (cs *cssStyle) surface(st css.Style) canvas.Color {
	if fill := cs.fill(st); fill.A() != 0 {
		return fill
	}
	return cs.win.Theme().Background
}
