package template

// Home collapses the stack back to the first screen that was shown, wherever
// it is. It is the "start over" of a navigation stack.
func (c *Context) Home() {
	if len(c.stack) > 1 {
		c.stack = c.stack[:1]
	}
}
