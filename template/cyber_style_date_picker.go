package template

func (cyberStyle) DatePicker(p DatePickerParams) {
	win := p.Win
	cv := win.Canvas()
	u := popupScale(p.W)
	v := vividFor(win)

	glow(cv, halo{x: p.X, y: p.Y, w: p.W, h: p.H, radius: 4 * u, colour: v.neon, strength: 12, u: u})
	body := cyberPalette.panel
	if p.State.Hovered {
		body = cyberPalette.panelHov
	}
	neonFill(cv, p.X, p.Y, p.W, p.H, 4*u, body, v.edge)

	dp := datePaint{x: p.X, y: p.Y, w: p.W, h: p.H, u: u, month: p.Month, year: p.Year,
		firstWD: p.FirstWD, days: p.Days, selected: p.Selected, today: p.Today, hover: p.Hover}
	paintCyberDateHeader(cv, v, dp)
	paintWeekdayHeads(cv, dp, v.dim)
	paintCyberDateDays(cv, v, dp)
}
