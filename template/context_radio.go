package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

func (c *Context) Radio(area canvas.Area, params antui.RadioParams) event.Event {
	return c.tpl.Radio(c.win, area, params)
}
