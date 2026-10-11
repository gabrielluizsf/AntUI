package template

import (
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

// Switch draws a flip toggle like a checkbox but a sliding thumb. It flips
// value when clicked and reports whether it flipped this frame.
func (c *CSS) Switch(label string, value *bool) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleSwitch, label)
	if !ok {
		return event.Nothing
	}
	u := c.style.u()
	cy := y + (h-20*u)/2
	on := value != nil && *value
	s := c.interactionState("switch", x, cy, 36*u, 20*u, label, 36*u+8*u+textWidth(u, label), 20*u, State{On: on})
	st = c.style.entryStyle(c.style.cur, css.RoleSwitch, s)
	if !interactive(st) {
		c.paint(st, css.RoleSwitch, label, x, y, w, h, func() {
			c.style.Switch(c.win, s, x, cy, label, on)
		})
		return event.Nothing
	}
	var e event.Event
	c.paint(st, css.RoleSwitch, label, x, y, w, h, func() {
		e = c.ui.Switch(c.win, x, cy, label, value)
	})
	return e
}
