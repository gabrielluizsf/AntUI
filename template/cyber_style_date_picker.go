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

	p := datePaint{x: x, y: y, w: w, h: h, u: u, month: month, year: year,
		firstWD: firstWD, days: days, selected: selected, today: today, hover: hover}
	paintCyberDateHeader(cv, v, p)
	paintWeekdayHeads(cv, p, v.dim)
	paintCyberDateDays(cv, v, p)
}
