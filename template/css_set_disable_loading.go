package template

// SetDisableLoading turns the start-up loading screen off or back on. When
// it goes off after the screen was already presented, the picture it left
// is wiped, so a program that asks to be without it does not sit on it
// until its first frame. See [Template.SetDisableLoading].
func (c *CSS) SetDisableLoading(disabled bool) {
	shown := c.ui.ld.active()
	c.ui.SetDisableLoading(disabled)
	if disabled && shown {
		c.ui.ld.wipe(c.win, c.ui)
	}
}
