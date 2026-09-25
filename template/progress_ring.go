package template

import (
	"math"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// progressRing paints the progress as a ring: a full track annulus from a
// darker fill (or an explicit border colour), with a clockwise accent arc of
// the value sweeping from the top of the ring. The ring's thickness follows
// border-width, or a slice of the box when no width is declared. The widget's
// label, if any, floats in the middle.
func (cs *cssStyle) progressRing(win *antui.Window, st css.Style, x, y, w, h int, value float32) {
	u := cs.u()
	cv := win.Canvas()
	cx, cy := x+w/2, y+h/2
	radius := min(w, h) / 2
	thickness := max(2, radius/8)
	if st.Has("border-width") {
		thickness = max(st.BorderWidth[0]*u, 1)
	}
	track := canvas.Shade(cs.fill(st), -0.25)
	if st.Has("border-color") && st.BoxColor[0].A() != 0 {
		track = canvas.Fade(st.BoxColor[0], cs.alpha(st))
	}
	cv.RingArc(cx, cy, radius, thickness, -math.Pi/2, 2*math.Pi, track)
	if value > 0 {
		cv.RingArc(cx, cy, radius, thickness, -math.Pi/2, 2*math.Pi*float64(value), cs.fill(st))
	}
}
