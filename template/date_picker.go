package template

import "github.com/gabrielluizsf/antui"

func (builtinStyle) DatePicker(win *antui.Window, s State, x, y, w, h int, year, month, firstWD, days, selected, today, hover int) {
	cv := win.Canvas()
	th := win.Theme()
	u := Scale(win)

	cv.FillRoundRect(x, y, w, h, 4*u, th.Surface)
	border := th.Border
	if s.Focused {
		border = th.Accent
	}
	cv.RoundRect(x, y, w, h, 4*u, border)

	paintDatePickerHeader(cv, th, u, x, y, w, dateHeaderHeight(u), month, year)
	paintWeekdayHeads(cv, u, x, y, firstWD, th.TextMuted)
	paintDateDays(cv, th, u, x, y, firstWD, days, selected, today, hover)
}
