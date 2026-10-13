package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) paintBoxOf(win *antui.Window, st css.Style, x, y, w, h int, surface bool) {
	cv := win.Canvas()
	rx := st.Radius[css.SideTop] * cs.u()
	ry := rx
	if st.RadiusY[css.SideTop] != 0 {
		ry = st.RadiusY[css.SideTop] * cs.u()
	}

	if len(st.BackdropFilters) > 0 {
		cs.applyFilters(cv, st.BackdropFilters, x, y, w, h)
	}

	cs.paintShadows(win, st, x, y, w, h, rx, ry, false)

	hasBg := surface && cs.backgroundOn(st)
	if hasBg {
		cs.paintBackground(win, st, x, y, w, h, rx, ry)
	}

	cs.paintShadows(win, st, x, y, w, h, rx, ry, true)

	if st.BorderOn() {
		cs.paintBorder(cv, st, x, y, w, h, rx, ry, cs.interior(st, hasBg))
	}

	if cs.outlineOn(st) {
		cs.paintOutline(win, st, x, y, w, h)
	}

	if len(st.Filters) > 0 {
		cs.applyFilters(cv, st.Filters, x, y, w, h)
	}
}
