package template

import "github.com/gabrielluizsf/antui"

func (cyberStyle) DatePicker(win *antui.Window, s State, x, y, w, h int, year, month, firstWD, days, selected, today, hover int) {
	cv := win.Canvas()
	u := Scale(win)
	v := vividFor(win)

	glow(cv, x, y, w, h, 4*u, v.neon, 12, u)
	body := cyberPalette.panel
	if s.Hovered {
		body = cyberPalette.panelHov
	}
	neonFill(cv, x, y, w, h, 4*u, body, v.edge)

	// Header: month and year, flanked by the arrows.
	header := dateHeaderHeight(u)
	mname := monthName(month)
	labelW := textWidth(u, mname) + 3*u + textWidth(u, itoa(year))
	lx := x + (w-labelW)/2
	cv.TextScaled(lx, y+(header-textHeight(u))/2, mname, v.text, u)
	cv.TextScaled(lx+textWidth(u, mname)+3*u, y+(header-textHeight(u))/2, itoa(year), v.dim, u)
	cv.Line(x+3*u, y+header/2, x+6*u, y+header/2, v.neon)
	cv.Line(x+2*u, y+header/2, x+5*u, y+header/2, v.neon)
	cv.Line(x+w-3*u, y+header/2, x+w-6*u, y+header/2, v.neon)
	cv.Line(x+w-2*u, y+header/2, x+w-5*u, y+header/2, v.neon)

	// Column heads.
	c := dateCellSize(u)
	gx, _, _, _ := dateGridRect(x, y, u)
	hy := weekdayHeadY(x, y, u)
	for d := 0; d < 7; d++ {
		cx := gx + d*c
		cv.TextScaled(cx+(c-textWidth(u, "S"))/2, hy, dayHeadName((firstWD+d)%7), v.dim, u)
	}

	// Day cells.
	for day := 1; day <= days; day++ {
		pos := firstWD + day - 1
		col, row := pos%7, pos/7
		cx, cy, cc, _ := dayCellRect(x, y, u, col, row)
		switch day {
		case selected:
			glow(cv, cx, cy, cc, cc, 2*u, v.neon, int(float64(20)*v.glow), u)
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, v.neon)
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), cyberPalette.bg, u)
		case hover:
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, cyberPalette.panelHov)
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), v.text, u)
		default:
			colour := v.text
			if day == today {
				colour = v.neon
			}
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), colour, u)
		}
		if day == today && day != selected {
			cv.Circle(cx+cc/2, cy+cc/2, cc/2-2*u, v.neon)
		}
	}
}
