package template

import "github.com/gabrielluizsf/antui/canvas"

import "github.com/gabrielluizsf/antui/template/event"

// Button, Checkbox, Radio, Slider and Input draw the template's components.
// Each reports an [event.Event]; screens use them to decide what happened.
func (c *Context) Button(area canvas.Area, label string) event.Event {
	return c.tpl.Button(c.win, area, label)
}
