package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

// DatePicker draws a date picker that works like a select: a box showing the
// chosen date, with a calendar that pops up below it when the box is clicked.
// The header arrows (or PageUp/PageDown) flip through the months, and a click
// or Enter picks the highlighted day. Picking sets value and closes the popup;
// Escape or a click outside cancels. It reports when the calendar opened,
// closed, or a new day was picked.
func (t *uiTemplate) DatePicker(win *antui.Window, area canvas.Area, value *Date) event.Event {
	x, y, w, h := area.X, area.Y, area.Width, area.Height
	if value == nil || w <= 0 || h <= 0 {
		return event.Nothing
	}
	normalizeDate(value)
	id := win.WidgetID("template:date", canvas.Area{X: x, Y: y, Width: w, Height: h}, "")
	hovered := win.Hovered(canvas.Area{X: x, Y: y, Width: w, Height: h})
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
	t.style.DatePickerBox(DatePickerBoxParams{
		Win:   win,
		State: State{Hovered: hovered, Focused: focused, On: v.open},
		X:     x, Y: y, W: w, H: h,
		Value: dateString(*value),
		Open:  v.open,
	})
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
