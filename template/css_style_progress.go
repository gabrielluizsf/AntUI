package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Progress(p ProgressParams) {
	st := cs.styleFor(css.RoleProgress, State{})
	value := clampFloat(p.Value, 0, 1)
	u := cs.u()

	// A square progress whose border radius rounds it into a full circle is
	// a ring: the value draws as an arc sweeping around a track, instead of
	// a bar filling from the left.
	if cs.ringProgress(css.RoleProgress, st, p.W, p.H) {
		cs.progressRing(p.Win, st, p.X, p.Y, p.W, p.H, value)
		return
	}

	cv := p.Win.Canvas()
	fill := cs.fill(st)
	r := st.Radius[css.SideTop] * u
	if r <= 0 {
		r = p.H / 2
	}
	filled := int(float32(p.W)*value + 0.5)
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, r), canvas.Shade(fill, trackShade))
	if filled > 0 {
		cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: filled, Height: p.H}, r), fill)
	}
}

// ringProgress reports whether a progress of this box draws as a ring rather
// than a bar: the box is a square whose radius rounds it into a full circle.
func (cs *cssStyle) ringProgress(role string, st css.Style, w, h int) bool {
	if role != css.RoleProgress || w != h {
		return false
	}
	return st.Radius[css.SideTop]*cs.u() >= min(w, h)/2
}
