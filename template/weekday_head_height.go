package template

// weekdayHeadHeight is the row of column heads above the day grid: as tall
// as one glyph, so the heads never spill into the cells.
func weekdayHeadHeight(u int) int { return 16 * u }
