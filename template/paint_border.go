package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintBorder strokes the border. A uniform, fully-solid border is drawn as a
// thicker rounded rectangle behind the surface, which follows the radius.
// Mixed or non-solid borders are painted side by side with square corners.
func (cs *cssStyle) paintBorder(cv *canvas.Canvas, st css.Style, x, y, w, h, rx, ry int, fill canvas.Color) {
	u := cs.u()

	if !cs.perSideBorder(st) && len(st.BackgroundImages) == 0 {
		bw := st.BorderWidth[css.SideTop] * u
		if bw <= 0 {
			bw = u
		}
		col := cs.fade(st, st.BoxColor[0])
		if fill.A() == 0 {
			// There is no surface to paint over the middle of the ring, so the
			// ring is stroked instead of filled and left covered: one outline
			// per unit of width, which follows the radius the way the fill did.
			for i := range bw {
				cv.RoundRectXY(x+i, y+i, w-2*i, h-2*i, max(rx-i, 0), max(ry-i, 0), col)
			}
			return
		}
		cv.FillRoundRectXY(x, y, w, h, max(rx, u), max(ry, u), col)
		cv.FillRoundRectXY(x+bw, y+bw, w-2*bw, h-2*bw, max(rx-u, 0), max(ry-u, 0), fill)
		return
	}
	for side := 0; side < 4; side++ {
		style := st.BoxStyle[side]
		if style == css.BorderNone {
			continue
		}
		bw := st.BorderWidth[side] * u
		if bw <= 0 {
			bw = u
		}
		cs.paintSide(cv, side, style, bw, x, y, w, h, cs.fade(st, st.BoxColor[side]))
	}
}
