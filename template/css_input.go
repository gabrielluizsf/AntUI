package template

import (
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

// Input draws a single-line text field and reports whether its text changed.
func (c *CSS) Input(text *string) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleInput, "")
	if !ok {
		return event.Nothing
	}
	s := c.interactionState("input", x, y, w, h, "", w, h, State{})
	st = c.style.entryStyle(c.style.cur, css.RoleInput, s)
	var e event.Event
	c.paint(st, css.RoleInput, "", x, y, w, h, func() {
		e = c.ui.Input(c.win, x, y, w, h, text)
	})
	return e
}
