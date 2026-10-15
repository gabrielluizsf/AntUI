package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

// Select draws a dropdown picker: a box showing the currently chosen option,
// and a menu of ALL options when it is open. Clicking the box opens or closes
// it; picking an option sets index to its position and closes the menu.
func (c *CSS) Select(index *int, options []string) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleSelect, "")
	if !ok {
		return event.Nothing
	}
	s := c.interactionState(widgetBox{"select", "", x, y, w, h, w, h}, State{})
	st = c.style.entryStyle(c.style.cur, css.RoleSelect, s)
	var e event.Event
	c.paint(st, css.RoleSelect, "", x, y, w, h, func() {
		e = c.ui.Select(c.win, canvas.Area{X: x, Y: y, Width: w, Height: h}, SelectParams{Index: index, Options: options})
	})
	return e
}
