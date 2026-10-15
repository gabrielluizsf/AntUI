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
	cv.TextScaled(canvas.Point{X: float64(x), Y: float64(y)}, mname, canvas.TextStyle{Color: th.Text, Scale: u})
	cv.TextScaled(canvas.Point{X: float64(x + textWidth(u, mname) + 4*u), Y: float64(y)}, itoa(p.year), canvas.TextStyle{Color: th.TextMuted, Scale: u})

	ax, ay, aw := dateArrowRect(p.x, p.y, u, true)
	drawIcon(cv, IconCaretLeft, ax, ay, aw, aw, th.TextMuted)
	ax, ay, aw = dateArrowRect(p.x, p.y, u, false)
	drawIcon(cv, IconCaretRight, ax, ay, aw, aw, th.TextMuted)
}
