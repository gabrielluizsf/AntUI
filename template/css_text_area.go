package template

import (
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

// TextArea draws a multi-line text field and reports whether its text
// changed. Enter inserts a line break; word wrap and scrolling follow.
func (c *CSS) TextArea(text *string) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleTextArea, "")
	if !ok {
		return event.Nothing
	}
	s := c.interactionState("textarea", x, y, w, h, "", w, h, State{})
	st = c.style.entryStyle(c.style.cur, css.RoleTextArea, s)
	var e event.Event
	c.paint(st, css.RoleTextArea, "", x, y, w, h, func() {
		e = c.ui.TextArea(c.win, x, y, w, h, text)
	})
	return e
}
