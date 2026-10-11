package template

// SetDisableLoading turns the loading screen off or back on, and when it
// turns it off after the screen was already presented, wipes the picture it
// left behind so the window shows the plain background instead. See
// [Template.SetDisableLoading].
func (c *Context) SetDisableLoading(disabled bool) {
	if c.tpl == nil {
		return
	}
	ld := c.tpl.loadingState()
	shown := ld.active()
	c.tpl.SetDisableLoading(disabled)
	if disabled && shown {
		ld.wipe(c.win, c.tpl)
	}
}
