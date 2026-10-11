package template

// Depth is how many screens are stacked, the current one included.
func (c *Context) Depth() int { return len(c.stack) }
