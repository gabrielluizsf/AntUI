package template

// shiftDate and shiftDateMonth move a single (year, month, day) by days or
// months, returning the new triple.
func shiftDate(year, month, day, delta int) (int, int, int) {
	d := Date{Year: year, Month: month, Day: day}
	shiftDay(&d, delta)
	return d.Year, d.Month, d.Day
}
