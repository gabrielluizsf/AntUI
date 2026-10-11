package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Progress(win *antui.Window, x, y, w, h int, value float32) {
	st := cs.styleFor(css.RoleProgress, State{})
	value = clampFloat(value, 0, 1)
	u := cs.u()

	// A square progress whose border radius rounds it into a full circle is
	// a ring: the value draws as an arc sweeping around a track, instead of
	// a bar filling from the left.
	r := st.Radius[0] * u
	if w == h && r >= min(w, h)/2 {
		cs.progressRing(win, st, x, y, w, h, value)
		return
	}

	cv := win.Canvas()
	fill := cs.fill(st)
	if r <= 0 {
		r = h / 2
	}
	filled := int(float32(w)*value + 0.5)
	cv.FillRoundRect(x, y, w, h, r, canvas.Shade(fill, -0.25))
	if filled > 0 {
		cv.FillRoundRect(x, y, filled, h, r, fill)
	}
}
