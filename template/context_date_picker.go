package template

import "github.com/gabrielluizsf/antui/template/event"

func (c *Context) DatePicker(x, y, w, h int, value *Date) event.Event {
	return c.tpl.DatePicker(c.win, x, y, w, h, value)
}
