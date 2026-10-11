package template

import "github.com/gabrielluizsf/antui/template/event"

func (c *Context) Radio(x, y int, label string, value *int, option int) event.Event {
	return c.tpl.Radio(c.win, x, y, label, value, option)
}
