package template

import "github.com/gabrielluizsf/antui/template/event"

func (c *Context) Input(x, y, w, h int, text *string) event.Event {
	return c.tpl.Input(c.win, x, y, w, h, text)
}
