package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

func (c *Context) Switch(area canvas.Area, params antui.CheckboxParams) event.Event {
	return c.tpl.Switch(c.win, area, params)
}
