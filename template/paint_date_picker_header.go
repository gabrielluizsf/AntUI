package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// paintDatePickerHeader draws the month and year, and the prev/next arrows in
// the gutters beside the day grid.
func paintDatePickerHeader(cv *canvas.Canvas, th *antui.Theme, p datePaint) {
	u := p.u
	mname := monthName(p.month)
	x := p.x + p.w/2 - textWidth(u, mname)/2
	y := p.headerTextY(textHeight(u))
	cv.TextScaled(x, y, mname, th.Text, u)
	cv.TextScaled(x+textWidth(u, mname)+4*u, y, itoa(p.year), th.TextMuted, u)

	ax, ay, aw := dateArrowRect(p.x, p.y, u, true)
	drawIcon(cv, IconCaretLeft, ax, ay, aw, aw, th.TextMuted)
	ax, ay, aw = dateArrowRect(p.x, p.y, u, false)
	drawIcon(cv, IconCaretRight, ax, ay, aw, aw, th.TextMuted)
}
