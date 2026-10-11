package template

import "github.com/gabrielluizsf/antui/template/event"

func (c *Context) TextArea(x, y, w, h int, text *string) event.Event {
	return c.tpl.TextArea(c.win, x, y, w, h, text)
}
