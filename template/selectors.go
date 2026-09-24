package template

import (
	"time"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

// Select draws a dropdown picker. The box shows the currently chosen option;
// clicking it opens or closes a menu of all options right below it. Picking an
// option sets index to its position. It reports when the menu opened, closed,
// or a new option was picked.
func (t *uiTemplate) Select(win *antui.Window, x, y, w, h int, index *int, options []string) event.Event {
	if index == nil || len(options) == 0 {
		return event.Nothing
	}
	if *index < 0 {
		*index = 0
	}
	if *index >= len(options) {
		*index = len(options) - 1
	}

	id := win.WidgetID("template:select", x, y, w, h, "")
	hovered := win.Hovered(x, y, w, h)
	open := t.selectOpen == id
	wasOpen := open
	menuH := h * len(options)

	// A press anywhere outside the box and its menu closes it. This runs
	// first so a click that closes the menu does not also count for the box.
	if open && win.MousePressed(antui.MouseLeft) && !win.Hovered(x, y, w, h+menuH) {
		open = false
	}

	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)

	picked := false

	// Keyboard, while the menu is open and the box holds the focus: arrows
	// walk the menu, Enter or Space picks it, Escape closes it. It runs before
	// the box's own click activation below, so Enter picks instead of just
	// toggling the menu.
	if open && focused {
		if win.KeyPressed(antui.KeyUp) && t.selectCursor > 0 {
			t.selectCursor--
		}
		if win.KeyPressed(antui.KeyDown) && t.selectCursor < len(options)-1 {
			t.selectCursor++
		}
		if win.KeyPressed(antui.KeyEnter) || win.KeyPressed(antui.KeySpace) {
			*index = t.selectCursor
			open, t.selectOpen = false, 0
			picked = true
		}
		if win.KeyPressed(antui.KeyEscape) {
			open, t.selectOpen = false, 0
		}
	}

	// Open/close with click or Enter/Space. Opening also grabs the keyboard,
	// so the arrows above work from the very first frame the menu is out.
	if wasOpen == open && !picked {
		if win.WidgetKeyActivate(id) || win.WidgetClick(id, hovered) {
			open = !open
			if open {
				t.selectCursor = *index
				t.selectOpen = id
				win.WidgetFocusSet(id)
			} else {
				t.selectOpen = 0
			}
		}
	}

	// Mouse: each menu entry is its own clickable area.
	if open {
		type optState struct {
			x, y  int
			opt   string
			sel   bool
			state State
		}
		opts := make([]optState, 0, len(options))
		for i, opt := range options {
			oy := y + h*(i+1)
			oid := win.WidgetID("template:option", x, oy, w, h, opt)
			oh := win.Hovered(x, oy, w, h)
			sel := i == *index
			if win.WidgetClick(oid, oh) {
				*index = i
				open, t.selectOpen = false, 0
				picked = true
			}
			st := State{Hovered: oh, Focused: focused, On: sel}
			if open && focused && t.selectCursor == i {
				st.Hovered = true
			}
			opts = append(opts, optState{x: x, y: oy, opt: opt, sel: sel, state: st})
		}
		// The menu is drawn dead last, above every widget the flow painted
		// under it, so an open menu never hides behind fields below.
		win.WidgetFrontmost(func() {
			for _, o := range opts {
				t.style.SelectOption(win, o.state, o.x, o.y, w, h, o.opt, o.sel)
			}
		})
	}

	t.style.Select(win, State{
		Hovered: hovered,
		Focused: focused,
		On:      open,
	}, x, y, w, h, options[*index], open)

	if picked {
		return t.fire(event.Event{Component: event.Select, Kind: event.Pick, Step: *index})
	}
	if wasOpen != open {
		kind := event.Close
		if open {
			kind = event.Open
		}
		return t.fire(event.Event{Component: event.Select, Kind: kind})
	}
	return event.Nothing
}

// Date is one day on the calendar the DatePicker shows. Month is 1..12, Day
// 1..31; the picker keeps the date valid as it moves around.
type Date struct {
	Year  int
	Month int
	Day   int
}

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

// dateString is the date a closed picker box shows.
func dateString(d Date) string {
	day := itoa(d.Day)
	if d.Day < 10 {
		day = "0" + day
	}
	return monthName(d.Month) + " " + day + ", " + itoa(d.Year)
}

// shiftDate and shiftDateMonth move a single (year, month, day) by days or
// months, returning the new triple.
func shiftDate(year, month, day, delta int) (int, int, int) {
	d := Date{Year: year, Month: month, Day: day}
	shiftDay(&d, delta)
	return d.Year, d.Month, d.Day
}

func shiftDateMonth(year, month, day, delta int) (int, int, int) {
	d := Date{Year: year, Month: month, Day: day}
	shiftMonth(&d, delta)
	return d.Year, d.Month, d.Day
}

// normalizeDate clamps a Date into a valid calendar day.
func normalizeDate(d *Date) {
	if d.Month < 1 {
		d.Month = 1
	}
	if d.Month > 12 {
		d.Month = 12
	}
	if d.Year == 0 {
		d.Year = 1970
	}
	if d.Day < 1 {
		d.Day = 1
	}
	if max := daysInMonth(d.Year, d.Month); d.Day > max {
		d.Day = max
	}
}

// shiftDay moves d by delta days, walking across month and year boundaries.
func shiftDay(d *Date, delta int) {
	d.Day += delta
	for d.Day < 1 {
		shiftMonth(d, -1)
		d.Day += daysInMonth(d.Year, d.Month)
	}
	for d.Day > daysInMonth(d.Year, d.Month) {
		d.Day -= daysInMonth(d.Year, d.Month)
		shiftMonth(d, 1)
	}
}

// shiftMonth moves d by delta months, clamping the day to the target month.
func shiftMonth(d *Date, delta int) {
	month := d.Month + delta
	for month < 1 {
		month += 12
		d.Year--
	}
	for month > 12 {
		month -= 12
		d.Year++
	}
	d.Month = month
	if max := daysInMonth(d.Year, d.Month); d.Day > max {
		d.Day = max
	}
}

// daysInMonth is the number of days in a calendar month.
func daysInMonth(year, month int) int {
	switch month {
	case 4, 6, 9, 11:
		return 30
	case 2:
		if leapYear(year) {
			return 29
		}
		return 28
	}
	return 31
}

func leapYear(y int) bool {
	return y%4 == 0 && (y%100 != 0 || y%400 == 0)
}

// firstWeekday is the weekday of the first day of the month: 0 = Sunday
// through 6 = Saturday.
func firstWeekday(year, month int) int {
	m, y := month, year
	if m < 3 {
		m += 12
		y--
	}
	d, c := y%100, y/100
	h := (1 + 13*(m+1)/5 + d + d/4 + c/4 + 5*c) % 7
	// h is Zeller: 0 = Saturday, so shift to Sunday = 0.
	return (h + 6) % 7
}

// datePickerSize is the full width and height of the calendar at scale u.
// The day grid is seven columns by six rows of cells sized by
// [dateCellSize], between a 20u header row and a 16u weekday row — the
// height of one glyph, so the column heads fit without touching the grid.
func datePickerSize(u int) (w, h int) {
	w = 14*u + 7*dateCellSize(u)
	h = dateHeaderHeight(u) + weekdayHeadHeight(u) + 6*dateCellSize(u)
	return
}

// dateHeaderHeight is where the month/year row sits, right below the header.
func dateHeaderHeight(u int) int { return 20 * u }

// weekdayHeadHeight is the row of column heads above the day grid: as tall
// as one glyph, so the heads never spill into the cells.
func weekdayHeadHeight(u int) int { return 16 * u }

// dateCellSize is the side of one day cell: a two-digit day number in the
// built-in face is 16u wide, so a cell needs at least that much plus a
// little bearing to keep the digits off its neighbours.
func dateCellSize(u int) int { return 18 * u }

// dateGridRect is the day grid's rectangle (column head row not included).
func dateGridRect(x, y, u int) (gx, gy, gw, gh int) {
	gx = x + 7*u
	gy = y + dateHeaderHeight(u) + weekdayHeadHeight(u)
	c := dateCellSize(u)
	return gx, gy, c * 7, c * 6
}

// weekdayHeadY is where a column head starts, its top edge in window
// coordinates. gy is the grid's top, so the band sits directly above it.
func weekdayHeadY(x, y, u int) int {
	_, gy, _, _ := dateGridRect(x, y, u)
	return gy - weekdayHeadHeight(u)
}

// dayCellRect is the rectangle of the day at the given grid offset.
func dayCellRect(x, y, u, col, row int) (cx, cy, cw, ch int) {
	c := dateCellSize(u)
	gx, gy, _, _ := dateGridRect(x, y, u)
	return gx + col*c, gy + row*c, c, c
}

// inDateGrid reports whether (mx,my) is inside the day grid.
func inDateGrid(x, y, u, mx, my int) bool {
	gx, gy, gw, gh := dateGridRect(x, y, u)
	return mx >= gx && mx < gx+gw && my >= gy && my < gy+gh
}

// inDateArrow reports whether (mx,my) is on the left (prev, true) or right
// (next, false) arrow of the header. The arrow is a square under the header
// text baseline, at either end of the row.
func inDateArrow(x, y, u, mx, my int, prev bool) bool {
	ax, aw := x+2*u, 8*u
	if !prev {
		w, _ := datePickerSize(u)
		ax = x + w - 2*u - aw
	}
	return mx >= ax && mx < ax+aw && my >= y && my < y+dateHeaderHeight(u)
}

// dayAt returns the day number under (mx,my), or 0 when none.
func dayAt(x, y, u, mx, my, year, month int) int {
	if !inDateGrid(x, y, u, mx, my) {
		return 0
	}
	c := dateCellSize(u)
	gx, gy, _, _ := dateGridRect(x, y, u)
	col := (mx - gx) / c
	row := (my - gy) / c
	days := daysInMonth(year, month)
	day := row*7 + col - firstWeekday(year, month) + 1
	if day < 1 || day > days {
		return 0
	}
	return day
}
