package svg

import "github.com/gabrielluizsf/antui/canvas"

// palette is the colours a paint resolves against as the drawing is walked: the
// colour a `currentColor` takes. It travels with the walk because the colour is
// only known where the drawing is painted, not where it was read.
type palette struct {
	current canvas.Color
}
