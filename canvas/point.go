package canvas

// Point is one vertex of a path, in whatever units the path was built in —
// user units, a viewBox's coordinates, a reference-pixel grid. A path is
// geometry; only the canvas it is drawn on knows what a unit costs.
type Point struct {
	X, Y float64
}
