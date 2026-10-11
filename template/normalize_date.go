package template

// normalizeDate clamps a Date into a valid calendar day.
func normalizeDate(d *Date) {
	if d.Month < 1 {
		d.Month = 1
	}
	if d.Month > 12 {
		d.Month = 12
	}
	if d.Year == 0 {
		d.Year = 1970
	}
	if d.Day < 1 {
		d.Day = 1
	}
	if max := daysInMonth(d.Year, d.Month); d.Day > max {
		d.Day = max
	}
}
