package template

// dateGridRect is the day grid's rectangle (column head row not included).
func dateGridRect(x, y, u int) (gx, gy, gw, gh int) {
	gx = x + 7*u
	gy = y + dateHeaderHeight(u) + weekdayHeadHeight(u)
	c := dateCellSize(u)
	return gx, gy, c * 7, c * 6
}
