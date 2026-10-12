package template

// dateView is what one frame of the picker shows: which month and day, whether
// the calendar is open, and whether this frame picked a day.
type dateView struct {
	year, month, day int
	open, changed    bool
}
