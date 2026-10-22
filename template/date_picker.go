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

	// A calendar is a popover above the rest of the page, which is the one
	// shadow in this theme with real reach: nothing else here floats free of
	// the surface it was drawn on.
	cv.BoxShadow(canvas.BoxShadowParams{X: p.X, Y: p.Y, Width: p.W, Height: p.H, RX: 4 * u, RY: 4 * u, DY: 4 * u, Blur: 16 * u, Color: th.Shadow})
	cv.FillRoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, 4*u), th.Surface)
	border := th.Border
	if p.State.Focused {
		border = th.Accent
	}
	cv.RoundRect(canvas.Round(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: p.H}, 4*u), border)

	// The month, the year and the arrows that flip between them are the
	// calendar's own toolbar, set apart from the grid of days the way a
	// sheet's header sits on a tier of its own rather than blending into
	// the page it opens onto.
	cv.FillRect(canvas.Area{X: p.X, Y: p.Y, Width: p.W, Height: dateHeaderHeight(u)}, th.SurfaceHover)

	dp := datePaint{x: p.X, y: p.Y, w: p.W, h: p.H, u: u, month: p.Month, year: p.Year,
		firstWD: p.FirstWD, days: p.Days, selected: p.Selected, today: p.Today, hover: p.Hover, pick: p.Pick}
	paintDatePickerHeader(cv, th, dp)
	if p.Pick {
		paintMonthGrid(cv, th, dp)
		return
	}
	paintWeekdayHeads(cv, dp, th.TextMuted)
	paintDateDays(cv, th, dp)
}
