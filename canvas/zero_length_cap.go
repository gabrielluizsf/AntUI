package canvas

// zeroLengthCap is what a subpath that came down to a single point covers. It is
// the cap drawn at that point with no direction to reach along, so a round one is
// a disc and a square one a box, while a butt cap leaves nothing to fill.
func zeroLengthCap(sub []Point, half float64, cap LineCap) []Point {
	if len(sub) != 1 {
		return nil
	}
	switch cap {
	case CapRound:
		return capDisc(sub[0], half)
	case CapSquare:
		return capOutline(sub[0], Point{1, 0}, half, cap)
	}
	return nil
}
