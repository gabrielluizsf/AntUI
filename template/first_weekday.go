package template

// firstWeekday is the weekday of the first day of the month: 0 = Sunday
// through 6 = Saturday.
func firstWeekday(year, month int) int {
	m, y := month, year
	if m < 3 {
		m += 12
		y--
	}
	d, c := y%100, y/100
	h := (1 + 13*(m+1)/5 + d + d/4 + c/4 + 5*c) % 7
	// h is Zeller: 0 = Saturday, so shift to Sunday = 0.
	return (h + 6) % 7
}
