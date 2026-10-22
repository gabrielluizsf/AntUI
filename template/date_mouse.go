package template

import "github.com/gabrielluizsf/antui/canvas"

import "github.com/gabrielluizsf/antui"

// dateMouse applies a press inside an open calendar: outside the box and
// popup closes it, the header arrows flip the month or — in the month grid a
// tap on the title opened — the year, the title itself toggles between the
// two grids, and a cell in whichever one is showing picks from it.
func (t *uiTemplate) dateMouse(win *antui.Window, value *Date, v dateView, g dateGeom, id uint32, focused bool) dateView {
	if !win.MousePressed(antui.MouseLeft) {
		return v
	}
	mx, my := win.MouseX(), win.MouseY()
	switch {
	case !win.Hovered(canvas.Area{X: g.x, Y: g.y, Width: g.w, Height: g.h}) && !win.Hovered(canvas.Area{X: g.px, Y: g.py, Width: g.dw, Height: g.dh}):
		v.open = false
	case inDateArrow(g.px, g.py, g.u, mx, my, true):
		if v.pickMonth {
			v.year--
		} else {
			v.year, v.month, v.day = shiftDateMonth(v.year, v.month, v.day, -1)
		}
	case inDateArrow(g.px, g.py, g.u, mx, my, false):
		if v.pickMonth {
			v.year++
		} else {
			v.year, v.month, v.day = shiftDateMonth(v.year, v.month, v.day, 1)
		}
	case inDateHeaderTitle(g.px, g.py, g.u, mx, my):
		v.pickMonth = !v.pickMonth
	case v.pickMonth:
		if m := monthAt(g.px, g.py, g.u, mx, my); m > 0 {
			v.month, v.pickMonth = m, false
			v.day = min(v.day, daysInMonth(v.year, v.month))
		}
	case inDateGrid(g.px, g.py, g.u, mx, my):
		if !focused {
			win.WidgetFocusSet(id)
		}
		if day := dayAt(g.px, g.py, g.u, mx, my, v.year, v.month); day > 0 {
			value.Year, value.Month, value.Day = v.year, v.month, day
			normalizeDate(value)
			v.changed, v.open = true, false
		}
	}
	return v
}
