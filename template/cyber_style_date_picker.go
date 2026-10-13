package template

import "github.com/gabrielluizsf/antui"

func (cyberStyle) DatePicker(win *antui.Window, s State, x, y, w, h int, year, month, firstWD, days, selected, today, hover int) {
	cv := win.Canvas()
	u := popupScale(w)
	v := vividFor(win)

	glow(cv, x, y, w, h, 4*u, v.neon, 12, u)
	body := cyberPalette.panel
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	neonFill(cv, x, y, w, h, 4*u, body, v.edge)

	paintCyberDateHeader(cv, v, u, x, y, w, dateHeaderHeight(u), month, year)
	paintWeekdayHeads(cv, u, x, y, firstWD, v.dim)
	paintCyberDateDays(cv, v, u, x, y, firstWD, days, selected, today, hover)
}
