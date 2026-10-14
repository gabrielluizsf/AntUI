package template

import (
	"github.com/gabrielluizsf/antui/canvas"
)

func (builtinStyle) DatePicker(p DatePickerParams) {
	win := p.Win
	cv := win.Canvas()
	th := win.Theme()
	// The popup is drawn at the scale it was measured at, which is the step it
	// grew to for this window rather than the window's own.
	u := popupScale(p.W)

	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, 4*u), th.Surface)
	border := th.Border
	if p.State.Focused {
		border = th.Accent
	}
	cv.RoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, 4*u), border)

	dp := datePaint{x: p.X, y: p.Y, w: p.W, h: p.H, u: u, month: p.Month, year: p.Year,
		firstWD: p.FirstWD, days: p.Days, selected: p.Selected, today: p.Today, hover: p.Hover}
	paintDatePickerHeader(cv, th, dp)
	paintWeekdayHeads(cv, dp, th.TextMuted)
	paintDateDays(cv, th, dp)
}
