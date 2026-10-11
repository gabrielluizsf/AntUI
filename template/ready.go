package template

// Ready tells the loading screen the program has finished loading; the next
// frame shows the program itself. See [Template.Ready].
func (c *Context) Ready() {
	if c.tpl != nil {
		c.tpl.Ready()
	}
}
