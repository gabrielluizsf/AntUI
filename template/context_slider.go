package template

import "github.com/gabrielluizsf/antui/template/event"

func (c *Context) Slider(x, y, w, h int, value *float32, minValue, maxValue float32) event.Event {
	return c.tpl.Slider(c.win, x, y, w, h, value, minValue, maxValue)
}
