package template

import "github.com/gabrielluizsf/antui/canvas"

func (cyberStyle) DatePicker(p DatePickerParams) {
	win := p.Win
	cv := win.Canvas()
	u := popupScale(p.W)
	v := vividFor(win)

	// The site's window frame: a thin accent wash, no external halo.
	body := cyberPalette.panel
	if p.State.Hovered {
		body = cyberPalette.panelHov
	}
	neonFill(cv, p.X, p.Y, p.W, p.H, 4*u, u, body, canvas.Fade(v.neon, 90))

	// A rule under the month and the arrows is the header's own edge — a
	// filled band would square off a corner the panel itself chamfers.
	headerY := p.Y + dateHeaderHeight(u)
	cv.Line(canvas.Point{X: float64(p.X), Y: float64(headerY)}, canvas.Point{X: float64(p.X + p.W), Y: float64(headerY)}, canvas.Fade(v.neon, 60))

	dp := datePaint{x: p.X, y: p.Y, w: p.W, h: p.H, u: u, month: p.Month, year: p.Year,
		firstWD: p.FirstWD, days: p.Days, selected: p.Selected, today: p.Today, hover: p.Hover, pick: p.Pick}
	paintCyberDateHeader(cv, v, dp)
	if p.Pick {
		paintCyberMonthGrid(cv, v, dp)
		return
	}
	paintWeekdayHeads(cv, dp, v.dim)
	paintCyberDateDays(cv, v, dp)
}
