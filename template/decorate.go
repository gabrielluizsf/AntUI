package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// decorate paints underline, line-through and overline across a run.
func (cs *cssStyle) decorate(cv *canvas.Canvas, st css.Style, x, y, tw, th, scale int) {
	if st.TextDecoration == 0 || tw <= 0 {
		return
	}
	ink := cs.ink(st)
	u := max(scale, 1)
	if st.TextDecoration&css.TextDecorationUnderline != 0 {
		cv.FillRect(x, y+th-2*u, tw, u, ink)
	}
	if st.TextDecoration&css.TextDecorationLineThrough != 0 {
		cv.FillRect(x, y+th/2, tw, u, ink)
	}
	if st.TextDecoration&css.TextDecorationOverline != 0 {
		cv.FillRect(x, y, tw, u, ink)
	}
}
