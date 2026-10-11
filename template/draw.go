package template

// Draw paints the current screen once, cleared to the template's background
// first. Call it every frame between Begin and End; [Context.Run] does that
// for the whole life of the window.
//
// The loading screen arms first, before anything the screen draws: its draw
// is registered ahead of any the screen registers, and Window.End runs them
// in reverse, so the screen lands above a dropdown or a popup the app put
// up while it was still loading. It arms even with no screen pushed at all —
// a program still loading has nothing to show, and the loading screen is
// the showing.
func (c *Context) Draw() {
	if c.win == nil {
		return
	}
	if c.tpl != nil {
		c.tpl.loadingState().arm(c.win, c.tpl)
	}
	if len(c.stack) == 0 {
		return
	}
	c.win.Canvas().Clear(c.tpl.Background(c.win))
	c.stack[len(c.stack)-1].Draw(c)
}
