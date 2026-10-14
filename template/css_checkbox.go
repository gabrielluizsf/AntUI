package template

import (
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

// Checkbox draws a labelled checkbox, flips value when clicked, and reports
// whether it flipped this frame.
func (c *CSS) Checkbox(label string, value *bool) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleCheckbox, label)
	if !ok {
		return event.Nothing
	}
	u := c.style.u()
	box := boxSize * u
	// The box and its label are one line standing in the middle of whatever
	// height the solver gave the cell, so a row of taller cells still has its
	// controls on one line.
	cy := y + (h-box)/2
	on := value != nil && *value
	s := c.interactionState(widgetBox{"checkbox", label, x, cy, box, box, box + labelGap*u + textWidth(u, label), box}, State{On: on})
	st = c.style.entryStyle(c.style.cur, css.RoleCheckbox, s)
	if !interactive(st) {
		c.paint(st, css.RoleCheckbox, label, x, y, w, h, func() {
			c.style.Checkbox(CheckboxParams{
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
	c.paint(st, css.RoleCheckbox, label, x, y, w, h, func() {
		e = c.ui.Checkbox(c.win, x, y, label, value)
	})
	return e
}
