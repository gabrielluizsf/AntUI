package template

// Current is the screen showing right now, or nil when none has been shown yet.
func (c *Context) Current() Screen {
	if len(c.stack) == 0 {
		return nil
	}
	return c.stack[len(c.stack)-1]
}
