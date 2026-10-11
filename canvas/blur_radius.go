package canvas

// blurRadius turns a CSS blur radius into the box radius three passes want.
// A Gaussian stands in for the CSS blur with a standard deviation of half the
// radius, but three box passes carry heavier shoulders than a Gaussian and
// bleed noticeably wider than a browser at the same σ. Taking a third of the
// radius brings the visible reach back to roughly a Gaussian of that σ, so a
// shadow hugs its box the way it does on the web.
func blurRadius(blur int) int {
	if blur <= 0 {
		return 0
	}
	return max(blur/3, 1)
}
