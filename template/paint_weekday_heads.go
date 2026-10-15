package template

import "github.com/gabrielluizsf/antui/canvas"

// paintWeekdayHeads draws the S M T W T F S column heads, each over the column
// it names — which starts part way into the week, so the first head may stand
// over an empty column when the month begins mid-week.
func paintWeekdayHeads(cv *canvas.Canvas, p datePaint, colour canvas.Color) {
	c := dateCellSize(p.u)
	gx, _, _, _ := dateGridRect(p.x, p.y, p.u)
	for d := 0; d < 7; d++ {
		cx := gx + d*c
		cv.TextScaled(canvas.Point{X: float64(cx + (c-textWidth(p.u, "S"))/2), Y: float64(p.weekHeadY())},
			dayHeadName((p.firstWD+d)%7), canvas.TextStyle{Color: colour, Scale: p.u})
	}
}
