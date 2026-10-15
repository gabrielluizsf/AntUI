package template

import "github.com/gabrielluizsf/antui/canvas"

import "github.com/gabrielluizsf/antui/template/event"

func (c *Context) DatePicker(area canvas.Area, value *Date) event.Event {
	return c.tpl.DatePicker(c.win, area, value)
}
