package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// paintMonthGrid draws the twelve months of the year a tap on the header's
// title opened, the chosen one filled the way a chosen day is.
func paintMonthGrid(cv *canvas.Canvas, th *antui.Theme, p datePaint) {
	inset := 3 * p.u
	for m := 1; m <= 12; m++ {
		col, row := (m-1)%4, (m-1)/4
		cx, cy, cw, ch := monthCellRect(p.x, p.y, p.u, col, row)
		bx, by, bw, bh := cx+inset, cy+inset, cw-2*inset, ch-2*inset
		name := monthAbbr(m)
		tx := bx + (bw-textWidth(p.u, name))/2
		ty := by + (bh-textHeight(p.u))/2
		if m == p.month {
			cv.FillRoundRect(canvas.Round(canvas.Area{X: bx, Y: by, Width: bw, Height: bh}, p.u*3), th.Accent)
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, name, canvas.TextStyle{Color: th.TextOnAccent, Scale: p.u})
		} else {
			cv.TextScaled(canvas.Point{X: float64(tx), Y: float64(ty)}, name, canvas.TextStyle{Color: th.Text, Scale: p.u})
		}
	}
}
