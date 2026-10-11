package template

import "github.com/gabrielluizsf/antui/template/event"

func (c *Context) Select(x, y, w, h int, index *int, options []string) event.Event {
	return c.tpl.Select(c.win, x, y, w, h, index, options)
}
