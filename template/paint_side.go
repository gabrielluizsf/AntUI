package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintSide strokes one side of a box in the border style the element declared,
// handing it to the painter for that style. The styles are separate painters
// rather than branches of one, so a new style is a new file and nothing else
// has to know about it.
func (cs *cssStyle) paintSide(cv *canvas.Canvas, s boxSide, style uint8) {
	switch style {
	case css.BorderDouble:
		cs.paintDoubleSide(cv, s)
	case css.BorderDashed:
		cs.paintDashSide(cv, s, false)
	case css.BorderDotted:
		cs.paintDashSide(cv, s, true)
	default:
		cs.paintSolidSide(cv, s)
	}
}
