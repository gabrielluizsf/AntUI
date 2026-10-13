package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

// DatePicker paints the calendar popup: the month header, the weekday row, and
// the day grid below it with the chosen, hovered and today's days marked.
func (cs *cssStyle) DatePicker(win *antui.Window, s State, x, y, w, h int, year, month, firstWD, days, selected, today, hover int) {
	// The calendar is a datepicker carrying a class of its own, so a sheet can
	// tell it apart from the box it dropped from while the element's rules go
	// on dressing them both.
	st := cs.styleForClass(css.RoleDatePicker, "popup", s)
	// The popup grew to a step of its own to fill the window, and its text has
	// to grow with the days it sits in: the scale below is what the calendar
	// measures and draws at, until this call is done.
	u := popupScale(w)
	cs.textScale = u
	defer func() { cs.textScale = 0 }()

	p := datePaint{x: x, y: y, w: w, h: h, u: u, month: month, year: year,
		firstWD: firstWD, days: days, selected: selected, today: today, hover: hover}
	cs.paintBox(win, st, x, y, w, h)
	cv := win.Canvas()
	cs.paintDateHeader(cv, st, p)
	cs.paintDateRules(cv, st, p)
	cs.paintDateDays(cv, st, p)
}
