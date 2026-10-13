package template

// dateArrowGutter is the column down each side of the calendar: the arrow, a
// little air outside it against the popup's edge, and the space between it and
// the day columns that keeps the arrows off the dates.
func dateArrowGutter(u int) int { return 24 * u }
