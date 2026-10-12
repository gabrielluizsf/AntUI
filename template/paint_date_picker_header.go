package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// paintDatePickerHeader draws the month and year between the prev/next arrows.
func paintDatePickerHeader(cv *canvas.Canvas, th *antui.Theme, u, x, y, w, header, month, year int) {
	mname := monthName(month)
	aw := 8 * u
	cv.TextScaled(x+w/2-textWidth(u, mname)/2, y+(header-textHeight(u))/2, mname, th.Text, u)
	cv.TextScaled(x+w/2-textWidth(u, mname)/2+textWidth(u, mname)+4*u, y+(header-textHeight(u))/2, itoa(year), th.TextMuted, u)

	// Arrow glyphs.
	prevX, nextX := x+2*u, x+w-2*u-aw
	ay := y + (header-textHeight(u))/2
	cv.Line(prevX+3*u, ay+textHeight(u)/2, prevX+aw-3*u, ay+3*u, th.TextMuted)
	cv.Line(prevX+3*u, ay+textHeight(u)/2, prevX+aw-3*u, ay+textHeight(u)-3*u, th.TextMuted)
	cv.Line(nextX+3*u, ay+3*u, nextX+aw-3*u, ay+textHeight(u)/2, th.TextMuted)
	cv.Line(nextX+aw-3*u, ay+textHeight(u)/2, nextX+3*u, ay+textHeight(u)-3*u, th.TextMuted)
}
