package template

import "time"

// SetLoadingTimeout changes how long the loading screen waits before the
// app closes. See [Template.SetLoadingTimeout].
func (c *Context) SetLoadingTimeout(d time.Duration) {
	if c.tpl != nil {
		c.tpl.SetLoadingTimeout(d)
	}
}
