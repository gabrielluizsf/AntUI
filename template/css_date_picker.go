package template

import (
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

// DatePicker draws a picker that works like a select: a box showing the chosen
// date, and a calendar that pops up below the box when it is clicked.
func (c *CSS) DatePicker(value *Date) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleDatePicker, "")
	if !ok {
		return event.Nothing
	}
	s := c.interactionState(widgetBox{"date", "", x, y, w, h, w, h}, State{})
	st = c.style.entryStyle(c.style.cur, css.RoleDatePicker, s)
	var e event.Event
	c.paint(st, css.RoleDatePicker, "", x, y, w, h, func() {
		e = c.ui.DatePicker(c.win, x, y, w, h, value)
	})
	return e
}
