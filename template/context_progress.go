package template

import "github.com/gabrielluizsf/antui/canvas"

func (c *Context) Progress(area canvas.Area, progress float32) {
	c.tpl.Progress(c.win, area, progress)
}
