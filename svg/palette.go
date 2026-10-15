package svg

import "github.com/gabrielluizsf/antui/canvas"

// palette is the colours a paint resolves against as the drawing is walked: the
// colour a `currentColor` takes, and what `context-fill` and `context-stroke`
// take inside the content of a marker — the fill and the stroke of the element
// that named the marker. The two context colours travel with the walk because
// the colour is only known where the drawing is painted, not where it was read,
// and they are empty everywhere but inside a marker's content, where there is
// no context element and a paint asking for one paints nothing, as the spec
// says.
type palette struct {
	current canvas.Color
	fill    canvas.Color
	stroke  canvas.Color
}
