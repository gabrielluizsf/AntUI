package canvas

// crossing is where a horizontal line meets an edge, and which way the edge was
// heading when it was met. Sorted by x, every second pair of crossings is a
// stretch of line that falls inside the shape.
type crossing struct {
	x   float64
	dir int
}
