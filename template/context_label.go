package template

// Label draws a line of text in the template's own voice.
func (c *Context) Label(x, y int, text string) { c.tpl.Label(c.win, x, y, text) }
