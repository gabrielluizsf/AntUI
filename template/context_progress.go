package template

func (c *Context) Progress(x, y, w, h int, progress float32) {
	c.tpl.Progress(c.win, x, y, w, h, progress)
}
