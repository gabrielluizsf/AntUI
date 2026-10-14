package template

import (
	"github.com/gabrielluizsf/antui/template/css"
)

// DatePicker paints the calendar popup: the month header, the weekday row, and
// the day grid below it with the chosen, hovered and today's days marked.
func (cs *cssStyle) DatePicker(p DatePickerParams) {
	// The calendar is a datepicker carrying a class of its own, so a sheet can
	// tell it apart from the box it dropped from while the element's rules go
	// on dressing them both.
	st := cs.styleForClass(css.RoleDatePicker, "popup", p.State)
	// The popup grew to a step of its own to fill the window, and its text has
	// to grow with the days it sits in: the scale below is what the calendar
	// measures and draws at, until this call is done.
	u := popupScale(p.W)
	cs.textScale = u
	defer func() { cs.textScale = 0 }()

	dp := datePaint{x: p.X, y: p.Y, w: p.W, h: p.H, u: u, month: p.Month, year: p.Year,
		firstWD: p.FirstWD, days: p.Days, selected: p.Selected, today: p.Today, hover: p.Hover}
	cs.paintBox(p.Win, st, p.X, p.Y, p.W, p.H)
	cv := p.Win.Canvas()
	cs.paintDateHeader(cv, st, dp)
	cs.paintDateRules(cv, st, dp)
	cs.paintDateDays(cv, st, dp)
}
