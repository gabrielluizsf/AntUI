package canvas

import "math"

// curveSegments is how many straight segments a curve of the given size is
// flattened into. The gap between a chord and the arc it stands for shrinks
// with the square of the number of segments, so the count follows the square
// root of the curve's size: a quarter turn twenty pixels across wants about
// eight, and a tiny one wants two — anything less and a corner that should be
// a curve turns into a visible facet.
func curveSegments(d float64) int {
	return max(int(math.Ceil(1.6*math.Sqrt(max(d, 0)))), 2)
}
