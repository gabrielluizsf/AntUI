package canvas

// makeEdge is a segment ready to be sampled: clipped to the rows the band is
// being filled over, its slope inverted into "how far right for one down", and
// the way it was heading recorded, which is what tells the nonzero rule which
// side is inside. A level edge is left out — a line lies along it for a whole
// row, and that is not a crossing but a tie.
//
// An edge lying wholly past one side of the band in x is still kept: a stretch
// of line it bounds may reach into the band from outside it, and without its
// crossing the pairs come out lopsided — the stretch that reaches in never
// gets a start, and the fill drops it rather than painting it, which empties
// the whole shape wherever a shape reaches past its band with a sloped edge.
// The crossing costs nothing even where the stretch it pairs up lies wholly
// outside: the coverage of a stretch is clipped to the band where it is added,
// so a pair that never comes near the band paints no pixel at all.
func makeEdge(a, b Point, y0, y1 int) (edge, bool) {
	if a.Y == b.Y {
		return edge{}, false
	}
	dir := 1
	if b.Y < a.Y {
		a, b, dir = b, a, -1
	}
	top, bottom := max(a.Y, float64(y0)), min(b.Y, float64(y1))
	if bottom <= top {
		return edge{}, false
	}
	e := edge{y0: top, y1: bottom, dxdy: (b.X - a.X) / (b.Y - a.Y), dir: dir}
	e.xAtY0 = a.X + (top-a.Y)*e.dxdy
	return e, true
}
