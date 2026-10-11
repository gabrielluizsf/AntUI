package template

// datePickerSize is the full width and height of the calendar at scale u.
// The day grid is seven columns by six rows of cells sized by
// [dateCellSize], between a 20u header row and a 16u weekday row — the
// height of one glyph, so the column heads fit without touching the grid.
func datePickerSize(u int) (w, h int) {
	w = 14*u + 7*dateCellSize(u)
	h = dateHeaderHeight(u) + weekdayHeadHeight(u) + 6*dateCellSize(u)
	return
}
