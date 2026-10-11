package template

// dayAt returns the day number under (mx,my), or 0 when none.
func dayAt(x, y, u, mx, my, year, month int) int {
	if !inDateGrid(x, y, u, mx, my) {
		return 0
	}
	c := dateCellSize(u)
	gx, gy, _, _ := dateGridRect(x, y, u)
	col := (mx - gx) / c
	row := (my - gy) / c
	days := daysInMonth(year, month)
	day := row*7 + col - firstWeekday(year, month) + 1
	if day < 1 || day > days {
		return 0
	}
	return day
}
