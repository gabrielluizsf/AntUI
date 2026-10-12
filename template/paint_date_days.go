package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// paintDateDays draws the day cells, marking the selected, hovered and today.
func paintDateDays(cv *canvas.Canvas, th *antui.Theme, u, x, y, firstWD, days, selected, today, hover int) {
	for day := 1; day <= days; day++ {
		pos := firstWD + day - 1
		col, row := pos%7, pos/7
		cx, cy, cc, _ := dayCellRect(x, y, u, col, row)
		switch day {
		case selected:
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, th.Accent)
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), th.TextOnAccent, u)
		case hover:
			cv.FillRoundRect(cx, cy, cc, cc, 2*u, th.SurfaceHover)
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), th.Text, u)
		default:
			colour := th.Text
			if day == today {
				colour = th.Accent
			}
			cv.TextScaled(cx+(cc-textWidth(u, itoa(day)))/2, cy+(cc-textHeight(u))/2, itoa(day), colour, u)
		}
		if day == today && day != selected {
			cv.Circle(cx+cc/2, cy+cc/2, cc/2, th.Accent)
		}
	}
}
