package template

import "time"

// dateToday is the day of the month that is today, or 0 when the shown month
// is not the current one.
func dateToday(year, month, days int) int {
	now := time.Now()
	if now.Year() == year && int(now.Month()) == month {
		return clampInt(now.Day(), 1, days)
	}
	return 0
}
