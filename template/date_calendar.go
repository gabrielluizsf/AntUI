package template

import "github.com/gabrielluizsf/antui"

// dateCalendar is what the popup needs to draw: the weekday it starts on, its
// length, the highlighted day clamped to it, and the day under the pointer.
func dateCalendar(win *antui.Window, v dateView, g dateGeom) (firstWD, days, sel, hover int) {
	firstWD = firstWeekday(v.year, v.month)
	days = daysInMonth(v.year, v.month)
	sel = clampInt(v.day, 1, days)
	if v.open && win.Hovered(g.px, g.py, g.dw, g.dh) {
		hover = dayAt(g.px, g.py, g.u, win.MouseX(), win.MouseY(), v.year, v.month)
	}
	return firstWD, days, sel, hover
}
