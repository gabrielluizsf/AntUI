package template

// Run is the whole application loop: it shows the first screen and keeps
// drawing it (and whatever it pushes) until the window closes — including
// the loading screen, which is what the window shows until Ready arrives
// and which closes the loop itself should the loading timeout pass first.
func (c *Context) Run(first Screen) {
	c.Go(first)
	for c.win.Begin() {
		c.Draw()
		c.win.End()
		if !c.win.Running() {
			break
		}
	}
}
