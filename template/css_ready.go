package template

// Ready tells the loading screen everything the program loads has loaded:
// the next frame is the program itself. Until then the loop keeps drawing
// the screen, and past the loading timeout the app closes. It may come from
// any goroutine. See [Template.Ready].
func (c *CSS) Ready() { c.ui.Ready() }
