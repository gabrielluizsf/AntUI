package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// paintDatePickerHeader draws the month and year, and the prev/next arrows in
// the gutters beside the day grid.
func paintDatePickerHeader(cv *canvas.Canvas, th *antui.Theme, u, x, y, w, header, month, year int) {
	mname := monthName(month)
	cv.TextScaled(x+w/2-textWidth(u, mname)/2, y+(header-textHeight(u))/2, mname, th.Text, u)
	cv.TextScaled(x+w/2-textWidth(u, mname)/2+textWidth(u, mname)+4*u, y+(header-textHeight(u))/2, itoa(year), th.TextMuted, u)

	ax, ay, aw := dateArrowRect(x, y, u, true)
	DrawIcon(cv, IconCaretLeft, ax, ay, aw, aw, th.TextMuted)
	ax, ay, aw = dateArrowRect(x, y, u, false)
	DrawIcon(cv, IconCaretRight, ax, ay, aw, aw, th.TextMuted)
}
