package template

import "github.com/gabrielluizsf/antui/canvas"

import "github.com/gabrielluizsf/antui/template/event"

func (c *Context) TextArea(area canvas.Area, text *string) event.Event {
	return c.tpl.TextArea(c.win, area, text)
}
