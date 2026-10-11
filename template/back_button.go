package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

// BackButton draws a small "Back" button in the top-left corner that pops the
// screen when clicked, and reports whether it was clicked. It is the usual
// way out of a screen that was pushed on top of another. It scales with the
// window like everything else.
func (c *Context) BackButton() event.Event {
	u := c.Scale()
	w := (24 + canvas.TextWidth("Back")) * u
	e := c.tpl.Button(c.win, 12*u, 12*u, w, 26*u, "Back")
	if e.Is(event.Button, event.Click) {
		c.Back()
	}
	return e
}
