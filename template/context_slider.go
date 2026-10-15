package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

func (c *Context) Slider(area canvas.Area, params antui.SliderParams) event.Event {
	return c.tpl.Slider(c.win, area, params)
}
