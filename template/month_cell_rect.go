package template

// monthCellRect is the rectangle of the month at the given grid offset in
// the month grid — four columns, three rows, the same bounds the day grid
// fills but divided into twelve cells instead of forty-two.
func monthCellRect(x, y, u, col, row int) (cx, cy, cw, ch int) {
	gx, gy, gw, gh := dateGridRect(x, y, u)
	cw, ch = gw/4, gh/3
	return gx + col*cw, gy + row*ch, cw, ch
}

// monthAt returns the month number (1-12) under (mx,my) in an open month
// grid, or 0 when none.
func monthAt(x, y, u, mx, my int) int {
	gx, gy, gw, gh := dateGridRect(x, y, u)
	if mx < gx || mx >= gx+gw || my < gy || my >= gy+gh {
		return 0
	}
	cw, ch := gw/4, gh/3
	col, row := (mx-gx)/cw, (my-gy)/ch
	m := row*4 + col + 1
	if m < 1 || m > 12 {
		return 0
	}
	return m
}
