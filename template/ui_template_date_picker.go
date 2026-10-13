package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

// DatePicker draws a date picker that works like a select: a box showing the
// chosen date, with a calendar that pops up below it when the box is clicked.
// The header arrows (or PageUp/PageDown) flip through the months, and a click
// or Enter picks the highlighted day. Picking sets value and closes the popup;
// Escape or a click outside cancels. It reports when the calendar opened,
// closed, or a new day was picked.
func (t *uiTemplate) DatePicker(win *antui.Window, x, y, w, h int, value *Date) event.Event {
	if value == nil || w <= 0 || h <= 0 {
		return event.Nothing
	}
	normalizeDate(value)
	id := win.WidgetID("template:date", x, y, w, h, "")
	hovered := win.Hovered(x, y, w, h)
	if hovered {
		win.WidgetHot(id)
	}
	g := dateGeomFor(win, x, y, w, h)
	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)

	v := dateView{year: value.Year, month: value.Month, day: value.Day, open: t.dateOpen == id}
	wasOpen := v.open
	v = t.dateInput(win, value, id, hovered, focused, g, v)

	firstWD, days, sel, hoverDay := dateCalendar(win, v, g)
	t.style.DatePickerBox(win, State{Hovered: hovered, Focused: focused, On: v.open},
		x, y, w, h, dateString(*value), v.open)
	if v.open {
		p := datePaint{x: g.px, y: g.py, w: g.dw, h: g.dh, u: g.u,
			year: v.year, month: v.month, firstWD: firstWD, days: days,
			selected: sel, today: dateToday(v.year, v.month, days), hover: hoverDay}
		t.showDatePopup(win, State{Focused: focused, Hovered: hovered}, p)
	}

	if v.changed {
		return t.fire(event.Event{Component: event.DatePicker, Kind: event.Change, Step: value.Day})
	}
	if wasOpen != v.open {
		kind := event.Close
		if v.open {
			kind = event.Open
		}
		return t.fire(event.Event{Component: event.DatePicker, Kind: kind})
	}
	return event.Nothing
}
