package template

import "github.com/gabrielluizsf/antui"

// dateKeys walks an open calendar with the keyboard: the arrows move the day a
// step, PageUp/PageDown flip a month, Home/End jump to the ends, Enter or
// Space picks the highlighted day and closes, Escape cancels.
func (t *uiTemplate) dateKeys(win *antui.Window, value *Date, v dateView) dateView {
	switch {
	case win.KeyPressed(antui.KeyLeft):
		v.year, v.month, v.day = shiftDate(v.year, v.month, v.day, -1)
	case win.KeyPressed(antui.KeyRight):
		v.year, v.month, v.day = shiftDate(v.year, v.month, v.day, 1)
	case win.KeyPressed(antui.KeyUp):
		v.year, v.month, v.day = shiftDate(v.year, v.month, v.day, -7)
	case win.KeyPressed(antui.KeyDown):
		v.year, v.month, v.day = shiftDate(v.year, v.month, v.day, 7)
	case win.KeyPressed(antui.KeyPageUp):
		v.year, v.month, v.day = shiftDateMonth(v.year, v.month, v.day, -1)
	case win.KeyPressed(antui.KeyPageDown):
		v.year, v.month, v.day = shiftDateMonth(v.year, v.month, v.day, 1)
	case win.KeyPressed(antui.KeyHome):
		v.day = 1
	case win.KeyPressed(antui.KeyEnd):
		v.day = daysInMonth(v.year, v.month)
	case win.KeyPressed(antui.KeyEnter), win.KeyPressed(antui.KeySpace):
		value.Year, value.Month, value.Day = v.year, v.month, v.day
		normalizeDate(value)
		v.changed, v.open = true, false
	case win.KeyPressed(antui.KeyEscape):
		v.open = false
	}
	return v
}
