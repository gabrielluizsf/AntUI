package template

// shiftDay moves d by delta days, walking across month and year boundaries.
func shiftDay(d *Date, delta int) {
	d.Day += delta
	for d.Day < 1 {
		shiftMonth(d, -1)
		d.Day += daysInMonth(d.Year, d.Month)
	}
	for d.Day > daysInMonth(d.Year, d.Month) {
		d.Day -= daysInMonth(d.Year, d.Month)
		shiftMonth(d, 1)
	}
}
