package template

import "github.com/gabrielluizsf/antui"

// showDatePopup draws the calendar dead last, above every widget the flow
// painted under it, so an open calendar never hides behind fields below.
func (t *uiTemplate) showDatePopup(win *antui.Window, state State, p datePaint) {
	win.WidgetFrontmost(func() {
		t.style.DatePicker(win, state, p.x, p.y, p.w, p.h,
			p.year, p.month, p.firstWD, p.days, p.selected, p.today, p.hover)
	})
}
