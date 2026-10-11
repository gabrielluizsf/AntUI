package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// maskLayer cuts an already-painted full-box layer down to the layer's
// background-clip region, keeping the box's corner curvature. A border-box
// clip masks the whole layer; a padding or content box moves the cut inside
// the border.
func (cs *cssStyle) maskLayer(layer *canvas.Canvas, st css.Style, clip uint8, x, y, w, h int, cr float64) {
	boxX, boxY, boxW, boxH := x, y, w, h
	switch clip {
	case css.BackPadding:
		boxX, boxY, boxW, boxH = cs.originRect(st, x, y, w, h, css.BackPadding)
	case css.BackContent:
		boxX, boxY, boxW, boxH = cs.originRect(st, x, y, w, h, css.BackContent)
	}
	if cr == 0 {
		layer.MaskRect(boxX, boxY, boxW, boxH)
		return
	}
	rr := min(boxW, boxH)
	layer.MaskRoundRect(boxX, boxY, boxW, boxH, int(cr*float64(rr)+0.5), int(cr*float64(rr)+0.5))
}
