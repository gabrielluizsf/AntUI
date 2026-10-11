package template

// weekdayHeadY is where a column head starts, its top edge in window
// coordinates. gy is the grid's top, so the band sits directly above it.
func weekdayHeadY(x, y, u int) int {
	_, gy, _, _ := dateGridRect(x, y, u)
	return gy - weekdayHeadHeight(u)
}
