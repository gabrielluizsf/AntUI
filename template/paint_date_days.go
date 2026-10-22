package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// paintDateDays draws the day cells, marking the selected, hovered and today.
// Each one is a circle a couple of units short of its cell on every side, not
// the cell's own square — the gap between circles is what tells two neighbours
// apart; a background tiled flush against the next one is what reads as a grid
// rather than a set of days.
func paintDateDays(cv *canvas.Canvas, th *antui.Theme, p datePaint) {
	inset := 2 * p.u
	for day := 1; day <= p.days; day++ {
		pos := p.firstWD + day - 1
		cx, cy, cc, _ := dayCellRect(p.x, p.y, p.u, pos%7, pos/7)
		ccx, ccy, r := cx+cc/2, cy+cc/2, cc/2-inset
		num := itoa(day)
		tx, ty := cx+(cc-textWidth(p.u, num))/2, cy+(cc-textHeight(p.u))/2
		switch day {
		case p.selected:
			cv.FillCircle(canvas.Circle{X: ccx, Y: ccy, R: r}, th.Accent)
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, num, canvas.TextStyle{Color: th.TextOnAccent, Scale: p.u})
		case p.hover:
			cv.FillCircle(canvas.Circle{X: ccx, Y: ccy, R: r}, th.SurfaceHover)
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, num, canvas.TextStyle{Color: th.Text, Scale: p.u})
		default:
			colour := th.Text
			if day == p.today {
				colour = th.Accent
			}
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, num, canvas.TextStyle{Color: colour, Scale: p.u})
		}
		if day == p.today && day != p.selected {
			cv.Circle(canvas.Circle{X: ccx, Y: ccy, R: r}, th.Accent)
		}
	}
}
