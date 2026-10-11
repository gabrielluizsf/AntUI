package template

import (
	"time"

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
	u := Scale(win)
	id := win.WidgetID("template:date", x, y, w, h, "")
	hovered := win.Hovered(x, y, w, h)
	if hovered {
		win.WidgetHot(id)
	}

	dw, dh := datePickerSize(u)
	px, py := x, y+h
	open := t.dateOpen == id
	wasOpen := open

	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)

	changed := false

	// The month an open calendar shows and the day it highlights; picking a
	// day copies them into value.
	curYear, curMonth, curDay := value.Year, value.Month, value.Day
	if open {
		curYear, curMonth, curDay = t.dateYear, t.dateMonth, t.dateCursor
	}

	// Keyboard, while the popup is open and the box holds focus: the arrows
	// walk the days (crossing a month boundary changes the shown month),
	// PageUp/PageDown flip months, Home/End jump to the first and last day,
	// Enter or Space picks the highlighted day and closes, Escape cancels.
	if open && focused {
		switch {
		case win.KeyPressed(antui.KeyLeft):
			curYear, curMonth, curDay = shiftDate(curYear, curMonth, curDay, -1)
		case win.KeyPressed(antui.KeyRight):
			curYear, curMonth, curDay = shiftDate(curYear, curMonth, curDay, 1)
		case win.KeyPressed(antui.KeyUp):
			curYear, curMonth, curDay = shiftDate(curYear, curMonth, curDay, -7)
		case win.KeyPressed(antui.KeyDown):
			curYear, curMonth, curDay = shiftDate(curYear, curMonth, curDay, 7)
		case win.KeyPressed(antui.KeyPageUp):
			curYear, curMonth, curDay = shiftDateMonth(curYear, curMonth, curDay, -1)
		case win.KeyPressed(antui.KeyPageDown):
			curYear, curMonth, curDay = shiftDateMonth(curYear, curMonth, curDay, 1)
		case win.KeyPressed(antui.KeyHome):
			curDay = 1
		case win.KeyPressed(antui.KeyEnd):
			curDay = daysInMonth(curYear, curMonth)
		case win.KeyPressed(antui.KeyEnter), win.KeyPressed(antui.KeySpace):
			value.Year, value.Month, value.Day = curYear, curMonth, curDay
			normalizeDate(value)
			changed = true
			open = false
		case win.KeyPressed(antui.KeyEscape):
			open = false
		}
	}

	// Mouse: a press outside the box and its open calendar closes it; the
	// header arrows flip the shown month; a day cell picks that day, sets
	// value and closes.
	if open && win.MousePressed(antui.MouseLeft) {
		mx, my := win.MouseX(), win.MouseY()
		switch {
		case !win.Hovered(x, y, w, h) && !win.Hovered(px, py, dw, dh):
			open = false
		case inDateArrow(px, py, u, mx, my, true):
			curYear, curMonth, curDay = shiftDateMonth(curYear, curMonth, curDay, -1)
		case inDateArrow(px, py, u, mx, my, false):
			curYear, curMonth, curDay = shiftDateMonth(curYear, curMonth, curDay, 1)
		case inDateGrid(px, py, u, mx, my):
			if !focused {
				win.WidgetFocusSet(id)
			}
			if day := dayAt(px, py, u, mx, my, curYear, curMonth); day > 0 {
				value.Year, value.Month, value.Day = curYear, curMonth, day
				normalizeDate(value)
				changed = true
				open = false
			}
		}
	}

	// The box: a click, or Enter/Space while it is closed, flips the popup.
	// While the popup is open a click on the box closes it again. Nothing the
	// popup already decided this frame is undone here.
	if wasOpen == open && !changed {
		if win.WidgetKeyActivate(id) || win.WidgetClick(id, hovered) {
			if open {
				open = false
			} else {
				curYear, curMonth, curDay = value.Year, value.Month, value.Day
				open = true
				win.WidgetFocusSet(id)
			}
		}
	}

	// While the popup is open, the pointer moving over the grid carries the
	// marker to whatever day it touches, so the picker shows exactly where a
	// click would land. The click above commits it. It only claims the marker
	// on frames the pointer actually moved, so the arrow keys keep control
	// while the mouse sits still.
	if open && win.Hovered(px, py, dw, dh) && win.MouseMoved() {
		if d := dayAt(px, py, u, win.MouseX(), win.MouseY(), curYear, curMonth); d > 0 {
			curDay = d
		}
	}

	if open {
		t.dateYear, t.dateMonth, t.dateCursor = curYear, curMonth, curDay
		t.dateOpen = id
	} else {
		t.dateOpen = 0
	}

	firstWD := firstWeekday(curYear, curMonth)
	days := daysInMonth(curYear, curMonth)
	sel := clampInt(curDay, 1, days)
	hoverDay := 0
	if open && win.Hovered(px, py, dw, dh) {
		hoverDay = dayAt(px, py, u, win.MouseX(), win.MouseY(), curYear, curMonth)
	}

	now := time.Now()
	today := 0
	if now.Year() == curYear && int(now.Month()) == curMonth {
		today = clampInt(now.Day(), 1, days)
	}

	t.style.DatePickerBox(win, State{Hovered: hovered, Focused: focused, On: open},
		x, y, w, h, dateString(*value), open)

	if open {
		// The popup is drawn dead last, above every widget the flow painted
		// under it, so an open calendar never hides behind fields below.
		state := State{Focused: focused, Hovered: hovered}
		yy, mm, wd, ds, sel, today, hover := curYear, curMonth, firstWD, days, sel, today, hoverDay
		win.WidgetFrontmost(func() {
			t.style.DatePicker(win, state, px, py, dw, dh, yy, mm, wd, ds, sel, today, hover)
		})
	}

	if changed {
		return t.fire(event.Event{Component: event.DatePicker, Kind: event.Change, Step: value.Day})
	}
	if wasOpen != open {
		kind := event.Close
		if open {
			kind = event.Open
		}
		return t.fire(event.Event{Component: event.DatePicker, Kind: kind})
	}
	return event.Nothing
}
