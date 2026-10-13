package template

// dateGridRect is the day grid's rectangle (column head row not included).
// The grid sits between the two arrow gutters, so it is centred in the box
// whatever the box is asked to be.
func dateGridRect(x, y, u int) (gx, gy, gw, gh int) {
	gx = x + dateArrowGutter(u)
	gy = y + dateHeaderHeight(u) + weekdayHeadHeight(u)
	c := dateCellSize(u)
	return gx, gy, c * 7, c * 6
}
