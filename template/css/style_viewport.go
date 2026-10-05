package css

// StyleViewport is [Sheet.Style] against a whole window: media queries test
// both of its edges — and its density and the color scheme its system paints
// in — while the viewport units resolve to the real size, so a stylesheet
// written for a tall narrow window sees that window instead of a square
// stand-in. The template computes every style of a frame against one
// viewport, which is what makes a window that resizes take effect on the
// next frame rather than halfway down the page.
func (sh *Sheet) StyleViewport(tag string, classes []string, state State, vp Viewport) Style {
	return sh.styleUnits(tag, classes, state, Units{
		Width: vp.Width, Height: vp.Height, Scale: vp.Scale, Scheme: vp.Scheme,
	}, true)
}
