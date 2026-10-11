package template

// Width and Height are the window's drawing size, for centering and for
// laying a screen out against its edges.
func (c *Context) Width() int { return c.win.Width() }
