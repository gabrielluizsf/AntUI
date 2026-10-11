package svg

// radius turns a length in pixels into the whole number of pixels a blur asks
// the canvas for.
func radius(v float64) int { return int(v + 0.5) }
