package template

import "github.com/gabrielluizsf/antui/canvas"

// paintWeekdayHeads draws the S M T W T F S column heads, fixed to the same
// seven columns every month — it is the day grid below that starts part way
// into the week, leaving the columns before day 1's own empty rather than
// the heads naming different days depending on where the month begins.
func paintWeekdayHeads(cv *canvas.Canvas, p datePaint, colour canvas.Color) {
	c := dateCellSize(p.u)
	gx, _, _, _ := dateGridRect(p.x, p.y, p.u)
	for d := 0; d < 7; d++ {
		cx := gx + d*c
		cv.TextScaled(canvas.Point{X: float64(cx + (c-textWidth(p.u, "S"))/2), Y: float64(p.weekHeadY())},
			dayHeadName(d), canvas.TextStyle{Color: colour, Scale: p.u})
	}
}
