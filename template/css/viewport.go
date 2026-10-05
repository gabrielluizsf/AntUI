package css

// Viewport is the window a stylesheet is read against: the size @media
// measures, the size the viewport units (vw/vh/vmin/vmax) resolve to, how
// many pixels the display draws for each point, and the color scheme the
// system paints in. The template latches one per frame, so every widget
// drawn in a frame is styled against the same window — a resize takes effect
// when the next frame opens instead of halfway down the page.
//
// Scale and Scheme are answers a system may not give: a zero Scale and a
// SchemeUnknown are "the display did not say", and the queries that read
// them are then treated as satisfied rather than guessed at.
type Viewport struct {
	Width, Height int

	// Scale is pixels per point — how many device pixels one CSS pixel of
	// this window covers — and 0 when the display did not say. It is what
	// the resolution media feature measures: 96 dots per inch at 1.
	Scale float64

	// Scheme is the color scheme prefers-color-scheme reads.
	Scheme Scheme
}

// resolutionDpi is the viewport's resolution in dots per inch, and whether
// the display gave one. A stylesheet's own numbers are rounded beside it, so
// a scale that lands on a repeating fraction (a display at 110 dpi is
// 1.145833… per point) still equals the number written for it.
func (vp Viewport) resolutionDpi() (float64, bool) {
	if vp.Scale <= 0 {
		return 0, false
	}
	return roundDpi(vp.Scale * 96), true
}

// roundDpi rounds dots per inch to the hundredth, which is finer than any
// display reports and far finer than the difference between two settings.
func roundDpi(dpi float64) float64 {
	return float64(int(dpi*100+0.5)) / 100
}
