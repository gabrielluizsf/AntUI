package template

// DatePickerParams holds parameters for painting a date picker calendar.
type DatePickerParams struct {
	WidgetGeometry
	Year, Month, FirstWD, Days, Selected, Today, Hover int
	Pick                                               bool // showing the month grid, not the day grid
}
