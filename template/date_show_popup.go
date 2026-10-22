package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// showDatePopup draws the calendar dead last, above every widget the flow
// painted under it, so an open calendar never hides behind fields below —
// behind a scrim over the whole window first, which is what tells the eye
// the calendar is the only thing asking for it right now.
func (t *uiTemplate) showDatePopup(win *antui.Window, state State, p datePaint) {
	win.WidgetFrontmost(func() {
		cv := win.Canvas()
		cv.FillRect(canvas.Area{X: 0, Y: 0, Width: win.Width(), Height: win.Height()}, canvas.Fade(canvas.Black, modalScrimAlpha))
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
			Pick:     p.pick,
		})
	})
}
