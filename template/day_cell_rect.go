package template

// dayCellRect is the rectangle of the day at the given grid offset.
func dayCellRect(x, y, u, col, row int) (cx, cy, cw, ch int) {
	c := dateCellSize(u)
	gx, gy, _, _ := dateGridRect(x, y, u)
	return gx + col*c, gy + row*c, c, c
}
