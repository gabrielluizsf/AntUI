package template

import "github.com/gabrielluizsf/antui"

// dateShowPopup draws the calendar dead last, above every widget the flow
// painted under it, so an open calendar never hides behind fields below.
func dateShowPopup(style Style, win *antui.Window, state State, g dateGeom, year, month, firstWD, days, sel, today, hover int) {
	win.WidgetFrontmost(func() {
		style.DatePicker(win, state, g.px, g.py, g.dw, g.dh, year, month, firstWD, days, sel, today, hover)
	})
}
