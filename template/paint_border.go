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

	if pts := clipPoints(st, x, y, w, h); pts != nil {
		// A clipped box is a polygon, so its border is the polygon's own
		// outline rather than a rectangle's radius.
		bw := st.BorderWidth[css.SideTop] * u
		if bw <= 0 {
			bw = u
		}
		strokeClip(cv, pts, bw, cs.fade(st, st.BoxColor[css.SideTop]))
		return
	}

	if !cs.perSideBorder(st) && len(st.BackgroundImages) == 0 {
		bw := st.BorderWidth[css.SideTop] * u
		if bw <= 0 {
			bw = u
		}
		col := cs.fade(st, st.BoxColor[css.SideTop])
		if fill.A() == 0 {
			// There is no surface to paint over the middle of the ring, so the
			// ring is stroked instead of filled and left covered: one outline
			// per unit of width, which follows the radius the way the fill did.
			for i := range bw {
				cv.RoundRect(canvas.Rounded{Area: canvas.Area{X: x + i, Y: y + i, Width: w - 2*i, Height: h - 2*i}, RX: max(rx-i, 0), RY: max(ry-i, 0)}, col)
			}
			return
		}
		cv.FillRoundRect(canvas.Rounded{Area: canvas.Area{X: x, Y: y, Width: w, Height: h}, RX: max(rx, u), RY: max(ry, u)}, col)
		cv.FillRoundRect(canvas.Rounded{Area: canvas.Area{X: x + bw, Y: y + bw, Width: w - 2*bw, Height: h - 2*bw}, RX: max(rx-u, 0), RY: max(ry-u, 0)}, fill)
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
		cs.paintSide(cv, boxSide{side: side, x: x, y: y, w: w, h: h,
			width: bw, colour: cs.fade(st, st.BoxColor[side])}, style)
	}
}
