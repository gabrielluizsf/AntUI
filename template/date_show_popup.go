package template

import "github.com/gabrielluizsf/antui"

// showDatePopup draws the calendar dead last, above every widget the flow
// painted under it, so an open calendar never hides behind fields below.
func (t *uiTemplate) showDatePopup(win *antui.Window, state State, p datePaint) {
	win.WidgetFrontmost(func() {
		t.style.DatePicker(DatePickerParams{
			Win:   win,
			State: state,
			X:     p.x, Y: p.y, W: p.w, H: p.h,
			Year:     p.year,
			Month:    p.month,
			FirstWD:  p.firstWD,
			Days:     p.days,
			Selected: p.selected,
			Today:    p.today,
			Hover:    p.hover,
		})
	})
}
