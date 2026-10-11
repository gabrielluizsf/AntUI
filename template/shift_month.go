package template

// shiftMonth moves d by delta months, clamping the day to the target month.
func shiftMonth(d *Date, delta int) {
	month := d.Month + delta
	for month < 1 {
		month += 12
		d.Year--
	}
	for month > 12 {
		month -= 12
		d.Year++
	}
	d.Month = month
	if max := daysInMonth(d.Year, d.Month); d.Day > max {
		d.Day = max
	}
}
