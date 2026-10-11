package svg

// offsets is how far along each side of the box the picture starts: 0 for the
// start edge, a half for the middle, a whole one for the far edge — the three
// positions each axis of the nine names has.
func (a PictureAlign) offsets() (fx, fy float64) {
	return float64(uint(a)%3) / 2, float64(uint(a)/3) / 2
}
