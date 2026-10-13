package template

import (
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

// Button draws a button and reports whether it was clicked.
func (c *CSS) Button(label string) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleButton, label)
	if !ok {
		return event.Nothing
	}
	s := c.interactionState(widgetBox{"button", label, x, y, w, h, w, h}, State{})
	st = c.style.entryStyle(c.style.cur, css.RoleButton, s)
	if !interactive(st) {
		c.paint(st, css.RoleButton, label, x, y, w, h, func() {
			c.style.Button(c.win, s, x, y, w, h, label)
		})
		return event.Nothing
	}
	var e event.Event
	c.paint(st, css.RoleButton, label, x, y, w, h, func() {
		e = c.ui.Button(c.win, x, y, w, h, label)
	})
	return e
}
