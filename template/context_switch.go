package template

import "github.com/gabrielluizsf/antui/template/event"

func (c *Context) Switch(x, y int, label string, value *bool) event.Event {
	return c.tpl.Switch(c.win, x, y, label, value)
}
