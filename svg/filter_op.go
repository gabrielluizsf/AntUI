package svg

import "github.com/gabrielluizsf/antui/canvas"

// filterOp is one step of the list the `filter` attribute holds: a function
// with what it was given — a fraction for grayscale, sepia, invert,
// brightness and contrast, degrees for hue-rotate, and a length in the
// drawing's own units for blur — or a reference to a `<filter>` element,
// which stands as the id it named until the whole drawing has been read and
// as the element itself after that. See [readFilters] and
// [Image.followFilterRefs].
type filterOp struct {
	kind   canvas.FilterKind
	amount float64
	ref    string
	def    *filterDef
	shadow *dropShadow
}
