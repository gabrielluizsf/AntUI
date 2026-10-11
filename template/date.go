package template

// Date is one day on the calendar the DatePicker shows. Month is 1..12, Day
// 1..31; the picker keeps the date valid as it moves around.
type Date struct {
	Year  int
	Month int
	Day   int
}
