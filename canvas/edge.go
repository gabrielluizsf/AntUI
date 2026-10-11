package canvas

// edge is one straight piece of a path, kept only where it can cross a line
// inside the band being filled.
type edge struct {
	y0, y1 float64
	dxdy   float64
	dir    int
	xAtY0  float64
}
