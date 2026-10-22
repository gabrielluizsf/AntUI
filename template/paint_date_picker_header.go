package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// paintDatePickerHeader draws the title — the month and year normally, just
// the year while the month grid a tap on it opened is showing, since the
// arrows beside it page a year at a time there instead of a month — and the
// prev/next arrows in the gutters beside the grid.
func paintDatePickerHeader(cv *canvas.Canvas, th *antui.Theme, p datePaint) {
	u := p.u
	y := p.headerTextY(textHeight(u))
	if p.pick {
		year := itoa(p.year)
		x := p.x + p.w/2 - textWidth(u, year)/2
		cv.TextScaled(canvas.Point{X: float64(x), Y: float64(y)}, year, canvas.TextStyle{Color: th.Text, Scale: u})
	} else {
		mname := monthName(p.month)
		x := p.x + p.w/2 - textWidth(u, mname)/2
		cv.TextScaled(canvas.Point{X: float64(x), Y: float64(y)}, mname, canvas.TextStyle{Color: th.Text, Scale: u})
		cv.TextScaled(canvas.Point{X: float64(x + textWidth(u, mname) + 4*u), Y: float64(y)}, itoa(p.year), canvas.TextStyle{Color: th.TextMuted, Scale: u})
	}

	ax, ay, aw := dateArrowRect(p.x, p.y, u, true)
	drawIcon(cv, IconCaretLeft, ax, ay, aw, aw, th.TextMuted)
	ax, ay, aw = dateArrowRect(p.x, p.y, u, false)
	drawIcon(cv, IconCaretRight, ax, ay, aw, aw, th.TextMuted)
}
