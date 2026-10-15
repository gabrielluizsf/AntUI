package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintDateHeader draws the title row: the month and year centred, the two
// arrows in the gutters, and a quiet band that separates navigation from the
// weekday and day rows below without introducing another card.
func (cs *cssStyle) paintDateHeader(cv *canvas.Canvas, st css.Style, p datePaint) {
	u, header := p.u, p.header()
	band := cs.fade(st, canvas.Shade(cs.fill(st), -0.035))
	line := cs.fade(st, canvas.Shade(cs.fill(st), -0.10))
	cv.FillRect(canvas.Area{X: p.x + u, Y: p.y + u, Width: p.w - 2*u, Height: max(header-u, 1)}, band)
	cv.FillRect(canvas.Area{X: p.x + 2*u, Y: p.y + header - u, Width: p.w - 4*u, Height: u}, line)

	mname := monthName(p.month)
	labelW := textWidth(u, mname) + 3*u + textWidth(u, itoa(p.year))
	lx := p.x + (p.w-labelW)/2
	ty := p.headerTextY(cs.textHeight(st))
	cs.run(cv, st, lx, ty, mname)
	cs.run(cv, st, lx+textWidth(u, mname)+3*u, ty, itoa(p.year))

	// The arrows sit in the header, leaving the weekday and day rows aligned as
	// one table below it.
	ax, ay, aw := dateArrowRect(p.x, p.y, u, true)
	drawIcon(cv, IconCaretLeft, ax, ay, aw, aw, cs.ink(st))
	ax, ay, aw = dateArrowRect(p.x, p.y, u, false)
	drawIcon(cv, IconCaretRight, ax, ay, aw, aw, cs.ink(st))
}
