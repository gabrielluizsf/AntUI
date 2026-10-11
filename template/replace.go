package template

// Replace swaps the top screen for another without going back first — the
// move a "logged in" screen makes after an entry form, so the form is not
// sitting underneath it waiting to be returned to.
func (c *Context) Replace(s Screen) {
	if s == nil {
		return
	}
	if len(c.stack) == 0 {
		c.stack = append(c.stack, s)
		return
	}
	c.stack[len(c.stack)-1] = s
}
