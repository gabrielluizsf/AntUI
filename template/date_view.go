package template

// dateView is what one frame of the picker shows: which month and day,
// whether the calendar is open, whether it is showing the day grid or the
// month grid a tap on the header switched it to, and whether this frame
// picked a day.
type dateView struct {
	year, month, day int
	open, changed    bool
	pickMonth        bool
}
