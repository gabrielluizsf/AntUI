package template

// Go puts a screen on top of the stack. It becomes the current screen on the
// next [Context.Draw].
func (c *Context) Go(s Screen) {
	if s != nil {
		c.stack = append(c.stack, s)
	}
}
