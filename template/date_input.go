package template

import "github.com/gabrielluizsf/antui"

// dateInput runs one frame of input for the picker — keyboard, mouse and the
// box toggle — and leaves the view it ends on, remembering it for the next
// frame.
func (t *uiTemplate) dateInput(win *antui.Window, value *Date, id uint32, hovered, focused bool, g dateGeom, v dateView) dateView {
	wasOpen := v.open
	if v.open {
		v.year, v.month, v.day, v.pickMonth = t.dateYear, t.dateMonth, t.dateCursor, t.datePickMonth
	}
	if v.open && focused && !v.pickMonth {
		v = t.dateKeys(win, value, v)
	}
	if v.open {
		v = t.dateMouse(win, value, v, g, id, focused)
	}
	if wasOpen == v.open && !v.changed {
		v = t.dateBoxClick(win, value, v, id, hovered)
	}
	if v.open && !v.pickMonth {
		v = t.dateHover(win, v, g)
	}
	if v.open {
		t.dateYear, t.dateMonth, t.dateCursor, t.dateOpen, t.datePickMonth = v.year, v.month, v.day, id, v.pickMonth
	} else {
		t.dateOpen = 0
	}
	return v
}
