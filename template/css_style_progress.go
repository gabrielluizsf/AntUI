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
	r := st.Radius[css.SideTop] * u
	if p.W == p.H && r >= min(p.W, p.H)/2 {
		cs.progressRing(p.Win, st, p.X, p.Y, p.W, p.H, value)
		return
	}

	cv := p.Win.Canvas()
	fill := cs.fill(st)
	if r <= 0 {
		r = p.H / 2
	}
	filled := int(float32(p.W)*value + 0.5)
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, r), canvas.Shade(fill, trackShade))
	if filled > 0 {
		cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: filled, Height: p.H}, r), fill)
	}
}
