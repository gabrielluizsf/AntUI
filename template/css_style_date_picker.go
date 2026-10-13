package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/css"
)

// datePaint is one calendar's geometry and its month, gathered so the painters
// that draw its parts agree on where they are and what day they are marking
// without passing thirteen numbers between them.
//
// x, y, w and h are the popup box, u the scale everything in it is measured at
// (its own, once it grew to fill the window), and the rest is the month: which
// weekday its first day falls on, how many days it has, and the day numbers to
// mark as chosen, today and hovered — 0 when not applicable.
type datePaint struct {
	x, y, w, h int
	u          int
	firstWD    int
	days       int
	selected   int
	today      int
	hover      int
}

// header is the height of the title row the month and the arrows sit in.
func (p datePaint) header() int { return dateHeaderHeight(p.u) }

// headerY is the top of the title row's text, centred in it.
func (p datePaint) headerTextY(th int) int { return p.y + (p.header()-th)/2 }

// weekHeadY is the top of the weekday-name row under the header.
func (p datePaint) weekHeadY() int { return weekdayHeadY(p.x, p.y, p.u) }

// rows is how many week rows the month actually needs, up to six.
func (p datePaint) rows() int { return min((p.firstWD+p.days+6)/7, 6) }

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

	p := datePaint{x: x, y: y, w: w, h: h, u: u, firstWD: firstWD, days: days,
		selected: selected, today: today, hover: hover}
	cs.paintBox(win, st, x, y, w, h)
	cv := win.Canvas()
	cs.paintDateHeader(cv, st, p, month, year)
	cs.paintDateRules(cv, st, p)
	cs.paintDateDays(cv, st, p)
}
