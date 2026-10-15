package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

func (c *Context) Select(area canvas.Area, params SelectParams) event.Event {
	return c.tpl.Select(c.win, area, params)
}
