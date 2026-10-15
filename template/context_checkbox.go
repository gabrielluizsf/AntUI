package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

func (c *Context) Checkbox(area canvas.Area, params antui.CheckboxParams) event.Event {
	return c.tpl.Checkbox(c.win, area, params)
}
