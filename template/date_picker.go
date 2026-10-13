package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) DatePicker(win *antui.Window, s State, x, y, w, h int, year, month, firstWD, days, selected, today, hover int) {
	cv := win.Canvas()
	th := win.Theme()
	// The popup is drawn at the scale it was measured at, which is the step it
	// grew to for this window rather than the window's own.
	u := popupScale(w)

	cv.FillRoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: w, Height: h}, 4*u), th.Surface)
	border := th.Border
	if s.Focused {
		border = th.Accent
	}
	cv.RoundRect(canvas.Round(canvas.Area{X: x, Y: y, Width: w, Height: h}, 4*u), border)

	p := datePaint{x: x, y: y, w: w, h: h, u: u, month: month, year: year,
		firstWD: firstWD, days: days, selected: selected, today: today, hover: hover}
	paintDatePickerHeader(cv, th, p)
	paintWeekdayHeads(cv, p, th.TextMuted)
	paintDateDays(cv, th, p)
}
