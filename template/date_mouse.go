package template

import "github.com/gabrielluizsf/antui"

// dateMouse applies a press inside an open calendar: outside the box and popup
// closes it, the header arrows flip months, a day cell picks that day.
func (t *uiTemplate) dateMouse(win *antui.Window, value *Date, v dateView, g dateGeom, id uint32, focused bool) dateView {
	if !win.MousePressed(antui.MouseLeft) {
		return v
	}
	mx, my := win.MouseX(), win.MouseY()
	switch {
	case !win.Hovered(g.x, g.y, g.w, g.h) && !win.Hovered(g.px, g.py, g.dw, g.dh):
		v.open = false
	case inDateArrow(g.px, g.py, g.u, mx, my, true):
		v.year, v.month, v.day = shiftDateMonth(v.year, v.month, v.day, -1)
	case inDateArrow(g.px, g.py, g.u, mx, my, false):
		v.year, v.month, v.day = shiftDateMonth(v.year, v.month, v.day, 1)
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
