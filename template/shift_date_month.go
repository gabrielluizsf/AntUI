package template

func shiftDateMonth(year, month, day, delta int) (int, int, int) {
	d := Date{Year: year, Month: month, Day: day}
	shiftMonth(&d, delta)
	return d.Year, d.Month, d.Day
}
