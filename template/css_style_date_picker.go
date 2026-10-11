package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) DatePicker(win *antui.Window, s State, x, y, w, h int, year, month, firstWD, days, selected, today, hover int) {
	st := cs.styleFor(css.RoleDatePicker, s)
	cs.paintBox(win, st, x, y, w, h)
	u := cs.u()

	header := dateHeaderHeight(u)
	mname := monthName(month)
	labelW := textWidth(u, mname) + 3*u + textWidth(u, itoa(year))
	lx := x + (w-labelW)/2
	cv := win.Canvas()
	cs.run(cv, st, lx, y+(header-cs.textHeight(st))/2, mname)
	cs.run(cv, st, lx+textWidth(u, mname)+3*u, y+(header-cs.textHeight(st))/2, itoa(year))

	aw := 8 * u
	ay := y + (header-cs.textHeight(st))/2
	prevX, nextX := x+2*u, x+w-2*u-aw
	cv.Line(prevX+3*u, ay+cs.textHeight(st)/2, prevX+aw-3*u, ay+3*u, cs.ink(st))
	cv.Line(prevX+3*u, ay+cs.textHeight(st)/2, prevX+aw-3*u, ay+cs.textHeight(st)-3*u, cs.ink(st))
	cv.Line(nextX+3*u, ay+3*u, nextX+aw-3*u, ay+cs.textHeight(st)/2, cs.ink(st))
	cv.Line(nextX+aw-3*u, ay+cs.textHeight(st)/2, nextX+3*u, ay+cs.textHeight(st)-3*u, cs.ink(st))

	c := dateCellSize(u)
	gx, _, _, _ := dateGridRect(x, y, u)
	hy := weekdayHeadY(x, y, u)
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
