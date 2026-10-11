package svg

import "github.com/gabrielluizsf/antui/canvas"

// markerVertex is one vertex of a path with the direction the path is going as
// it reaches and leaves it, which is everything `orient="auto"` needs to turn a
// marker to. hasIn and hasOut say whether there is a segment on that side at
// all: the first vertex of a path has nothing coming into it, the last has
// nothing going out, and the first point of a subpath after a move is not
// arrived at — the path does not travel from the end of one subpath to the
// start of the next.
type markerVertex struct {
	pt            canvas.Point
	hasIn, hasOut bool
	inX, inY      float64
	outX, outY    float64
}
