package template

// Back pops the top screen off and shows the one underneath. The first screen
// ever shown stays: Back does nothing there.
func (c *Context) Back() {
	if len(c.stack) > 1 {
		c.stack = c.stack[:len(c.stack)-1]
	}
}
