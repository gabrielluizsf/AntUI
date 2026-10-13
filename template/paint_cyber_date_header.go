package template

import "github.com/gabrielluizsf/antui/canvas"

// paintCyberDateHeader draws the month and year, and the neon arrows in the
// gutters beside the day grid.
func paintCyberDateHeader(cv *canvas.Canvas, v vivid, p datePaint) {
	u := p.u
	mname := monthName(p.month)
	labelW := textWidth(u, mname) + 3*u + textWidth(u, itoa(p.year))
	lx := p.x + (p.w-labelW)/2
	y := p.headerTextY(textHeight(u))
	cv.TextScaled(lx, y, mname, v.text, u)
	cv.TextScaled(lx+textWidth(u, mname)+3*u, y, itoa(p.year), v.dim, u)

	ax, ay, aw := dateArrowRect(p.x, p.y, u, true)
	cv.Line(ax+6*u, ay+2*u, ax+3*u, ay+aw/2, v.neon)
	cv.Line(ax+3*u, ay+aw/2, ax+6*u, ay+aw-2*u, v.neon)
	ax, ay, aw = dateArrowRect(p.x, p.y, u, false)
	cv.Line(ax+2*u, ay+2*u, ax+5*u, ay+aw/2, v.neon)
	cv.Line(ax+5*u, ay+aw/2, ax+2*u, ay+aw-2*u, v.neon)
}
