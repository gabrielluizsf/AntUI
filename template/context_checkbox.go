package template

import "github.com/gabrielluizsf/antui/template/event"

func (c *Context) Checkbox(x, y int, label string, value *bool) event.Event {
	return c.tpl.Checkbox(c.win, x, y, label, value)
}
