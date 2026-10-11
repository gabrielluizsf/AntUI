package template

import "github.com/gabrielluizsf/antui"

func (builtinStyle) DatePicker(win *antui.Window, s State, x, y, w, h int, year, month, firstWD, days, selected, today, hover int) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)

	cv.FillRoundRect(x, y, w, h, 4*u, th.Surface)
	border := th.Border
	if s.Focused {
		border = th.Accent
	}
	cv.RoundRect(x, y, w, h, 4*u, border)

	// Header: the prev/next arrows and the month name.
	mname := monthName(month)
	header := dateHeaderHeight(u)
	aw := 8 * u
	cv.TextScaled(x+w/2-textWidth(u, mname)/2, y+(header-textHeight(u))/2, mname, th.Text, u)
	cv.TextScaled(x+w/2-textWidth(u, mname)/2+textWidth(u, mname)+4*u, y+(header-textHeight(u))/2,
		itoa(year), th.TextMuted, u)

	// Arrow glyphs.
	prevX, nextX := x+2*u, x+w-2*u-aw
	ay := y + (header-textHeight(u))/2
	cv.Line(prevX+3*u, ay+textHeight(u)/2, prevX+aw-3*u, ay+3*u, th.TextMuted)
	cv.Line(prevX+3*u, ay+textHeight(u)/2, prevX+aw-3*u, ay+textHeight(u)-3*u, th.TextMuted)
	cv.Line(nextX+3*u, ay+3*u, nextX+aw-3*u, ay+textHeight(u)/2, th.TextMuted)
	cv.Line(nextX+aw-3*u, ay+textHeight(u)/2, nextX+3*u, ay+textHeight(u)-3*u, th.TextMuted)

	// Column heads S M T W T F S.
	c := dateCellSize(u)
	gx, _, _, _ := dateGridRect(x, y, u)
	hy := weekdayHeadY(x, y, u)
	for d := 0; d < 7; d++ {
		cx := gx + d*c
		cv.TextScaled(cx+(c-textWidth(u, "S"))/2, hy, dayHeadName((firstWD+d)%7), th.TextMuted, u)
	}

	// Day cells.
	for day := 1; day <= days; day++ {
		pos := firstWD + day - 1
		col, row := pos%7, pos/7
		cx, cy, cc, _ := dayCellRect(x, y, u, col, row)
		switch day {
		case selected:
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, th.Accent)
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), th.TextOnAccent, u)
		case hover:
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, th.SurfaceHover)
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), th.Text, u)
		default:
			colour := th.Text
			if day == today {
				colour = th.Accent
			}
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), colour, u)
		}
		if day == today && day != selected {
			cv.Circle(cx+cc/2, cy+cc/2, cc/2, th.Accent)
		}
	}
}
