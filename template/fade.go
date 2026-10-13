package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// fade is a colour the widget worked out for itself, faded by the style's
// opacity. A colour with no alpha of its own stays that way: fading it would
// hand the canvas the opaque black its colour carries underneath, painting
// there what the sheet left unpainted.
func (cs *cssStyle) fade(st css.Style, c canvas.Color) canvas.Color {
	if c.A() == 0 {
		return c
	}
	return fadeColor(c, cs.alpha(st))
}
