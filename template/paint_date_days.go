package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// paintDateDays draws the day cells, marking the selected, hovered and today.
func paintDateDays(cv *canvas.Canvas, th *antui.Theme, p datePaint) {
	for day := 1; day <= p.days; day++ {
		pos := p.firstWD + day - 1
		cx, cy, cc, _ := dayCellRect(p.x, p.y, p.u, pos%7, pos/7)
		num := itoa(day)
		tx, ty := cx+(cc-textWidth(p.u, num))/2, cy+(cc-textHeight(p.u))/2
		switch day {
		case p.selected:
			cv.FillRoundRect(cx, cy, cc, cc, 2*p.u, th.Accent)
			cv.TextScaled(tx, ty, num, th.TextOnAccent, p.u)
		case p.hover:
			cv.FillRoundRect(cx, cy, cc, cc, 2*p.u, th.SurfaceHover)
			cv.TextScaled(tx, ty, num, th.Text, p.u)
		default:
			colour := th.Text
			if day == p.today {
				colour = th.Accent
			}
			cv.TextScaled(tx, ty, num, colour, p.u)
		}
		if day == p.today && day != p.selected {
			cv.Circle(cx+cc/2, cy+cc/2, cc/2, th.Accent)
		}
	}
}
