package template

import "github.com/gabrielluizsf/antui/canvas"

// paintCyberDateHeader draws the month and year between the neon arrows.
func paintCyberDateHeader(cv *canvas.Canvas, v vivid, u, x, y, w, header, month, year int) {
	mname := monthName(month)
	labelW := textWidth(u, mname) + 3*u + textWidth(u, itoa(year))
	lx := x + (w-labelW)/2
	cv.TextScaled(lx, y+(header-textHeight(u))/2, mname, v.text, u)
	cv.TextScaled(lx+textWidth(u, mname)+3*u, y+(header-textHeight(u))/2, itoa(year), v.dim, u)
	cv.Line(x+3*u, y+header/2, x+6*u, y+header/2, v.neon)
	cv.Line(x+2*u, y+header/2, x+5*u, y+header/2, v.neon)
	cv.Line(x+w-3*u, y+header/2, x+w-6*u, y+header/2, v.neon)
	cv.Line(x+w-2*u, y+header/2, x+w-5*u, y+header/2, v.neon)
}
