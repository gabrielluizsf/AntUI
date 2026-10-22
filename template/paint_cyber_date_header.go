package template

import "github.com/gabrielluizsf/antui/canvas"

// paintCyberDateHeader draws the title — the month and year normally, just
// the year while the month grid a tap on it opened is showing, since the
// arrows beside it page a year at a time there instead of a month — and the
// neon arrows in the gutters beside the grid.
func paintCyberDateHeader(cv *canvas.Canvas, v vivid, p datePaint) {
	u := p.u
	y := p.headerTextY(textHeight(u))
	if p.pick {
		year := itoa(p.year)
		lx := p.x + (p.w-textWidth(u, year))/2
		cv.TextScaled(canvas.Point{X: float64(lx), Y: float64(y)}, year, canvas.TextStyle{Color: v.text, Scale: u})
	} else {
		mname := monthName(p.month)
		labelW := textWidth(u, mname) + 3*u + textWidth(u, itoa(p.year))
		lx := p.x + (p.w-labelW)/2
		cv.TextScaled(canvas.Point{X: float64(lx), Y: float64(y)}, mname, canvas.TextStyle{Color: v.text, Scale: u})
		cv.TextScaled(canvas.Point{X: float64(lx + textWidth(u, mname) + 3*u), Y: float64(y)}, itoa(p.year), canvas.TextStyle{Color: v.dim, Scale: u})
	}

	ax, ay, aw := dateArrowRect(p.x, p.y, u, true)
	drawIcon(cv, IconCaretLeft, ax, ay, aw, aw, v.neon)
	ax, ay, aw = dateArrowRect(p.x, p.y, u, false)
	drawIcon(cv, IconCaretRight, ax, ay, aw, aw, v.neon)
}
