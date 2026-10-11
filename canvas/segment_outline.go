package canvas

// segmentOutline is the rectangle of half-width on either side of one straight
// piece, ordered so that it always winds the same way round whatever direction
// the line runs: the normal side first, then across, then back.
func segmentOutline(a, b Point, half float64) []Point {
	d := unit(b.X-a.X, b.Y-a.Y)
	nx, ny := -d.Y*half, d.X*half
	return []Point{
		{a.X + nx, a.Y + ny},
		{b.X + nx, b.Y + ny},
		{b.X - nx, b.Y - ny},
		{a.X - nx, a.Y - ny},
	}
}
