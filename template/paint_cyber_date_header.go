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
	cv.TextScaled(canvas.Point{X: float64(lx), Y: float64(y)}, mname, canvas.TextStyle{Color: v.text, Scale: u})
	cv.TextScaled(canvas.Point{X: float64(lx + textWidth(u, mname) + 3*u), Y: float64(y)}, itoa(p.year), canvas.TextStyle{Color: v.dim, Scale: u})

	ax, ay, aw := dateArrowRect(p.x, p.y, u, true)
	cv.Line(canvas.Point{X: float64(ax + 6*u), Y: float64(ay + 2*u)}, canvas.Point{X: float64(ax + 3*u), Y: float64(ay + aw/2)}, v.neon)
	cv.Line(canvas.Point{X: float64(ax + 3*u), Y: float64(ay + aw/2)}, canvas.Point{X: float64(ax + 6*u), Y: float64(ay + aw - 2*u)}, v.neon)
	ax, ay, aw = dateArrowRect(p.x, p.y, u, false)
	cv.Line(canvas.Point{X: float64(ax + 2*u), Y: float64(ay + 2*u)}, canvas.Point{X: float64(ax + 5*u), Y: float64(ay + aw/2)}, v.neon)
	cv.Line(canvas.Point{X: float64(ax + 5*u), Y: float64(ay + aw/2)}, canvas.Point{X: float64(ax + 2*u), Y: float64(ay + aw - 2*u)}, v.neon)
}
