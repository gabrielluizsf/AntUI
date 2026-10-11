package template

import "time"

// SetLoadingTimeout changes how long the loading screen waits before the
// app closes, which is also the span its bar fills over. See
// [Template.SetLoadingTimeout].
func (c *CSS) SetLoadingTimeout(d time.Duration) { c.ui.SetLoadingTimeout(d) }
