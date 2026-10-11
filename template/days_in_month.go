package template

// daysInMonth is the number of days in a calendar month.
func daysInMonth(year, month int) int {
	switch month {
	case 4, 6, 9, 11:
		return 30
	case 2:
		if leapYear(year) {
			return 29
		}
		return 28
	}
	return 31
}
