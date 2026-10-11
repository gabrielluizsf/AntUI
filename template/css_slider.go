package template

import (
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

// Slider draws a horizontal slider and reports whether its value changed.
func (c *CSS) Slider(value *float32, minValue, maxValue float32) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleSlider, "")
	if !ok {
		return event.Nothing
	}
	s := c.interactionState("slider", x, y, w, h, "", w, h, State{})
	st = c.style.entryStyle(c.style.cur, css.RoleSlider, s)
	var e event.Event
	c.paint(st, css.RoleSlider, "", x, y, w, h, func() {
		e = c.ui.Slider(c.win, x, y, w, h, value, minValue, maxValue)
	})
	return e
}
