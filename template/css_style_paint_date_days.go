package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// paintDateDays draws one cell per day: the chosen day as a chip of the widget's
// own ink with the number reversed out in the surface colour, the hovered day
// as a quiet shade, and today as a dot in the middle underneath whichever of
// the two it happens to be.
func (cs *cssStyle) paintDateDays(cv *canvas.Canvas, st css.Style, p datePaint) {
	for day := 1; day <= p.days; day++ {
		c := dateCellAt(p, day, textWidth(p.u, itoa(day)), cs.textHeight(st))
		cs.paintDateDay(cv, st, p, day, c)
		if day == p.today && day != p.selected {
			cv.Circle(c.x+c.w/2, c.y+c.w/2, c.w/2-2*p.u, cs.ink(st))
		}
	}
}

// paintDateDay draws a single day's cell: its background by state, then the
// number on top in whichever colour that state calls for.
func (cs *cssStyle) paintDateDay(cv *canvas.Canvas, st css.Style, p datePaint, day int, c dateCell) {
	num := itoa(day)
	switch day {
	case p.selected:
		// The chosen day is a chip of the widget's own ink with the number
		// reversed out in the surface colour, the way the built-in templates
		// mark it in their accent colour.
		cv.FillRoundRect(c.x, c.y, c.w, c.h, 2*p.u, canvas.Fade(cs.ink(st), cs.alpha(st)))
		cs.runOn(cv, st, cs.surface(st), c.tx, c.ty, num)
	case p.hover:
		cv.FillRoundRect(c.x, c.y, c.w, c.h, 2*p.u, cs.fade(st, canvas.Shade(cs.fill(st), hoverFillShade)))
		cs.run(cv, st, c.tx, c.ty, num)
	default:
		cs.run(cv, st, c.tx, c.ty, num)
	}
}
