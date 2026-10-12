package svg

// boxBounds accumulates the union of the boxes it is given, in the drawing's
// own coordinates.
type boxBounds struct {
	minMax [4]float64
	found  bool
}

// add takes in one more box, growing the union to hold it.
func (b *boxBounds) add(x0, y0, x1, y1 float64) {
	if !b.found {
		b.minMax = [4]float64{x0, y0, x1, y1}
		b.found = true
		return
	}
	b.minMax[0] = min(b.minMax[0], x0)
	b.minMax[1] = min(b.minMax[1], y0)
	b.minMax[2] = max(b.minMax[2], x1)
	b.minMax[3] = max(b.minMax[3], y1)
}
