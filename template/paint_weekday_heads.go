package template

import "github.com/gabrielluizsf/antui/canvas"

// paintWeekdayHeads draws the S M T W T F S column heads.
func paintWeekdayHeads(cv *canvas.Canvas, u, x, y, firstWD int, colour canvas.Color) {
	c := dateCellSize(u)
	gx, _, _, _ := dateGridRect(x, y, u)
	hy := weekdayHeadY(x, y, u)
	for d := 0; d < 7; d++ {
		cx := gx + d*c
		cv.TextScaled(cx+(c-textWidth(u, "S"))/2, hy, dayHeadName((firstWD+d)%7), colour, u)
	}
}
