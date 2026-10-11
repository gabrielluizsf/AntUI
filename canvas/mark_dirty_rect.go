package canvas

// markDirtyRect widens the dirty bounds to cover a whole rectangle. The two
// opposite corners span it: the bounds math of markDirty makes any interior
// point redundant.
func (cv *Canvas) markDirtyRect(r Area) {
	if r.Width <= 0 || r.Height <= 0 {
		return
	}
	cv.markDirty(r.X, r.Y)
	cv.markDirty(r.X+r.Width-1, r.Y+r.Height-1)
}
