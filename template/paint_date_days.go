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
			cv.FillRoundRect(canvas.Round(canvas.Area{X: cx, Y: cy, Width: cc, Height: cc}, 2*p.u), th.Accent)
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, num, canvas.TextStyle{Color: th.TextOnAccent, Scale: p.u})
		case p.hover:
			cv.FillRoundRect(canvas.Round(canvas.Area{X: cx, Y: cy, Width: cc, Height: cc}, 2*p.u), th.SurfaceHover)
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, num, canvas.TextStyle{Color: th.Text, Scale: p.u})
		default:
			colour := th.Text
			if day == p.today {
				colour = th.Accent
			}
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, num, canvas.TextStyle{Color: colour, Scale: p.u})
		}
		if day == p.today && day != p.selected {
			cv.Circle(canvas.Circle{X: cx + cc/2, Y: cy + cc/2, R: cc / 2}, th.Accent)
		}
	}
}
