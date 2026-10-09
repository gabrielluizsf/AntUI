package css

// GetStyle computes the winning style for one widget in a window that is as
// tall as it is wide — [CSSClasses.GetStyleViewport] is the form that takes a
// real window, which is what a template drawing a frame uses. The result is
// cached, so a caller that paints every widget each frame computes each style
// once.
func (c *CSSClasses) GetStyle(tag string, classes []string, state State, width int) Style {
	return c.GetStyleViewport(tag, classes, state, Viewport{Width: width, Height: width})
}
