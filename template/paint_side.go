package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintSide strokes one side of a box in the given border style. side is 0
// top, 1 right, 2 bottom, 3 left.
func (cs *cssStyle) paintSide(cv *canvas.Canvas, side int, style uint8, bw, x, y, w, h int, col canvas.Color) {
	switch style {
	case css.BorderDouble:
		cs.paintDoubleSide(cv, side, bw, x, y, w, h, col)
	case css.BorderDashed:
		cs.paintDashSide(cv, side, bw, x, y, w, h, col, false)
	case css.BorderDotted:
		cs.paintDashSide(cv, side, bw, x, y, w, h, col, true)
	default:
		cs.paintSolidSide(cv, side, bw, x, y, w, h, col)
	}
}
