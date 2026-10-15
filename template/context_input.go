package template

import "github.com/gabrielluizsf/antui/canvas"

import "github.com/gabrielluizsf/antui/template/event"

func (c *Context) Input(area canvas.Area, text *string) event.Event {
	return c.tpl.Input(c.win, area, text)
}
