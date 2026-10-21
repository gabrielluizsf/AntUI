package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

// Radio draws one option of a group. Clicking an active option clears the
// group by storing antui.NoOption; otherwise it selects option.
func (c *CSS) Radio(label string, value *int, option int) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleRadio, label)
	if !ok {
		return event.Nothing
	}
	u := c.style.u()
	size := boxSize * u
	cy := y + (h-size)/2
	on := value != nil && *value == option
	s := c.interactionState(widgetBox{"radio", label, x, cy, size, option, size + labelGap*u + textWidth(u, label), size}, State{On: on})
	st = c.style.entryStyle(c.style.cur, css.RoleRadio, s)
	if !interactive(st) {
		c.paint(st, css.RoleRadio, label, x, y, w, h, func() {
			c.style.Radio(RadioParams{
				Win:   c.win,
				State: s,
				X:     x, Y: cy,
				Label: label,
				On:    on,
			})
		})
		return event.Nothing
	}
	var e event.Event
	c.paint(st, css.RoleRadio, label, x, y, w, h, func() {
		e = c.ui.Radio(c.win, canvas.Area{X: x, Y: cy}, antui.RadioParams{Label: label, Value: value, Option: option})
	})
	return e
}
