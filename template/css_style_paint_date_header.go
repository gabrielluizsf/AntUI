package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintDateHeader draws the title row: the month and year centred, the two
// arrows in the gutters, and a quiet band that separates navigation from the
// weekday and day rows below without introducing another card.
func (cs *cssStyle) paintDateHeader(cv *canvas.Canvas, st css.Style, p datePaint, month, year int) {
	u, header := p.u, p.header()
	band := cs.fade(st, canvas.Shade(cs.fill(st), -0.035))
	line := cs.fade(st, canvas.Shade(cs.fill(st), -0.10))
	cv.FillRect(p.x+u, p.y+u, p.w-2*u, max(header-u, 1), band)
	cv.FillRect(p.x+2*u, p.y+header-u, p.w-4*u, u, line)

	mname := monthName(month)
	labelW := textWidth(u, mname) + 3*u + textWidth(u, itoa(year))
	lx := p.x + (p.w-labelW)/2
	ty := p.headerTextY(cs.textHeight(st))
	cs.run(cv, st, lx, ty, mname)
	cs.run(cv, st, lx+textWidth(u, mname)+3*u, ty, itoa(year))

	// The arrows sit in the header, leaving the weekday and day rows aligned as
	// one table below it.
	ax, ay, aw := dateArrowRect(p.x, p.y, u, true)
	DrawIcon(cv, IconCaretLeft, ax, ay, aw, aw, cs.ink(st))
	ax, ay, aw = dateArrowRect(p.x, p.y, u, false)
	DrawIcon(cv, IconCaretRight, ax, ay, aw, aw, cs.ink(st))
}

// paintDateRules draws the fine grid the body is read as a table by: its rules
// make the weekday columns and week rows easy to scan while staying underneath
// the selected and hovered days, and the weekday names sit on top of them.
func (cs *cssStyle) paintDateRules(cv *canvas.Canvas, st css.Style, p datePaint) {
	u := p.u
	gx, gy, gw, _ := dateGridRect(p.x, p.y, u)
	hy := p.weekHeadY()
	c := dateCellSize(u)
	rule := cs.fade(st, canvas.Shade(cs.fill(st), -0.035))
	for col := 0; col <= 7; col++ {
		cv.FillRect(gx+col*c, hy, u, weekdayHeadHeight(u)+p.rows()*c, rule)
	}
	for row := 0; row <= p.rows(); row++ {
		cv.FillRect(gx, gy+row*c, gw, u, rule)
	}
	for d := 0; d < 7; d++ {
		cs.run(cv, st, gx+d*c+(c-textWidth(u, "S"))/2, hy, dayHeadName((p.firstWD+d)%7))
	}
}
