package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) DatePicker(win *antui.Window, s State, x, y, w, h int, year, month, firstWD, days, selected, today, hover int) {
	// The calendar is a datepicker carrying a class of its own, so a sheet can
	// tell it apart from the box it dropped from while the element's rules go
	// on dressing them both.
	st := cs.styleForClass(css.RoleDatePicker, "popup", s)
	// The popup grew to a step of its own to fill the window, and its text has
	// to grow with the days it sits in: the scale below is what the calendar
	// measures and draws at, until this call is done.
	u := popupScale(w)
	cs.textScale = u
	defer func() { cs.textScale = 0 }()
	cs.paintBox(win, st, x, y, w, h)

	header := dateHeaderHeight(u)
	// The title row is the table's heading: a quiet band separates navigation
	// from the weekday and day rows without introducing another card.
	band := cs.fade(st, canvas.Shade(cs.fill(st), -0.035))
	line := cs.fade(st, canvas.Shade(cs.fill(st), -0.10))
	cv := win.Canvas()
	cv.FillRect(x+u, y+u, w-2*u, max(header-u, 1), band)
	cv.FillRect(x+2*u, y+header-u, w-4*u, u, line)
	mname := monthName(month)
	labelW := textWidth(u, mname) + 3*u + textWidth(u, itoa(year))
	lx := x + (w-labelW)/2
	cs.run(cv, st, lx, y+(header-cs.textHeight(st))/2, mname)
	cs.run(cv, st, lx+textWidth(u, mname)+3*u, y+(header-cs.textHeight(st))/2, itoa(year))

	// The arrows sit in the header, leaving the weekday and day rows aligned as
	// one table below it.
	ax, ay, aw := dateArrowRect(x, y, u, true)
	DrawIcon(cv, IconCaretLeft, ax, ay, aw, aw, cs.ink(st))
	ax, ay, aw = dateArrowRect(x, y, u, false)
	DrawIcon(cv, IconCaretRight, ax, ay, aw, aw, cs.ink(st))

	c := dateCellSize(u)
	gx, gy, gw, _ := dateGridRect(x, y, u)
	hy := weekdayHeadY(x, y, u)
	// The calendar body is a table: its fine rules make the weekday columns
	// and week rows easy to scan while staying underneath selected/hovered days.
	rule := cs.fade(st, canvas.Shade(cs.fill(st), -0.035))
	rows := min((firstWD+days+6)/7, 6)
	for col := 0; col <= 7; col++ {
		cv.FillRect(gx+col*c, hy, u, weekdayHeadHeight(u)+rows*c, rule)
	}
	for row := 0; row <= rows; row++ {
		cv.FillRect(gx, gy+row*c, gw, u, rule)
	}
	for d := 0; d < 7; d++ {
		cx := gx + d*c
		cs.run(cv, st, cx+(c-textWidth(u, "S"))/2, hy, dayHeadName((firstWD+d)%7))
	}

	for day := 1; day <= days; day++ {
		pos := firstWD + day - 1
		col, row := pos%7, pos/7
		cx, cy, cc, _ := dayCellRect(x, y, u, col, row)
		tx, ty := cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-cs.textHeight(st))/2
		switch day {
		case selected:
			// The selected day is a chip of the widget's own ink with the
			// number reversed out in the surface colour, the way the
			// built-in templates mark it in their accent colour.
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, canvas.Fade(cs.ink(st), cs.alpha(st)))
			cs.runOn(cv, st, cs.surface(st), tx, ty, itoa(day))
		case hover:
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, cs.fade(st, canvas.Shade(cs.fill(st), -0.15)))
			cs.run(cv, st, tx, ty, itoa(day))
		default:
			cs.run(cv, st, tx, ty, itoa(day))
		}
		if day == today && day != selected {
			cv.Circle(cx+cc/2, cy+cc/2, cc/2-2*u, cs.ink(st))
		}
	}
}
