package template

// dateString is the date a closed picker box shows.
func dateString(d Date) string {
	day := itoa(d.Day)
	if d.Day < 10 {
		day = "0" + day
	}
	return monthName(d.Month) + " " + day + ", " + itoa(d.Year)
}
