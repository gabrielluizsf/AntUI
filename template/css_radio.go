package template

import (
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

// Radio draws one option of a group. Clicking it sets value to option and
// reports true; clicking the already-selected option changes nothing.
func (c *CSS) Radio(label string, value *int, option int) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleRadio, label)
	if !ok {
		return event.Nothing
	}
	u := c.style.u()
	size := boxSize * u
	cy := y + (h-size)/2
	on := value != nil && *value == option
	s := c.interactionState("radio", x, cy, size, option, label, size+labelGap*u+textWidth(u, label), size, State{On: on})
	st = c.style.entryStyle(c.style.cur, css.RoleRadio, s)
	if !interactive(st) {
		c.paint(st, css.RoleRadio, label, x, y, w, h, func() {
			c.style.Radio(c.win, s, x, cy, label, on)
		})
		return event.Nothing
	}
	var e event.Event
	c.paint(st, css.RoleRadio, label, x, y, w, h, func() {
		e = c.ui.Radio(c.win, x, cy, label, value, option)
	})
	return e
}
