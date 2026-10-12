package template

import "github.com/gabrielluizsf/antui"

// dateBoxClick flips the calendar when the box is clicked, or activated with
// Enter or Space.
func (t *uiTemplate) dateBoxClick(win *antui.Window, value *Date, v dateView, id uint32, hovered bool) dateView {
	if !win.WidgetKeyActivate(id) && !win.WidgetClick(id, hovered) {
		return v
	}
	if v.open {
		v.open = false
		return v
	}
	v.year, v.month, v.day = value.Year, value.Month, value.Day
	v.open = true
	win.WidgetFocusSet(id)
	return v
}
