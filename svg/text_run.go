package svg

import "github.com/gabrielluizsf/antui/canvas"

// TextRun is one piece of writing and everything that says where it goes: the
// pen it is written at, how far it is moved from where the piece before it left
// off, and the style it is written in. A run is what a `<tspan>` turns into, and
// the writing of a `<text>` with no `<tspan>` in it at all is a single run.
type TextRun struct {
	// Text is what is written.
	Text string
	// Style is how it is written: how tall, where along the pen it hangs, the
	// colour it is filled in, and the transform that turns the drawing's own
	// coordinates into the ones it lands on.
	Style Style
	// X and Y are where the pen goes before this piece, and HasX and HasY say
	// whether the drawing put it there: a run with no `x` of its own starts
	// where the writing around it left off, while `x="0"` is a place that has
	// been said and means the origin.
	X, Y float64
	// HasX and HasY tell an `x="0"` apart from no `x` at all, which is the
	// difference between writing at the edge of the drawing and writing on
	// from the last letter.
	HasX, HasY bool
	// XList and YList are a place for each letter rather than one place for
	// the whole of the piece, which is what `x="10 14 19"` says: the first
	// entry is where the first letter goes, the second where the second goes,
	// and a letter the list ran out before is left to follow the one before it.
	// They are only there when the drawing named more than one place — one
	// number is the X and Y above, which is what carries a single place across
	// the whitespace the file wrapped the writing with.
	XList, YList []float64
	// DX and DY are how far the pen is moved along from where it already was,
	// which is what spaces a piece away from the one before it without saying
	// anything about where on the page either of them is. They add up through
	// the elements around a piece, the way SVG reads them.
	DX, DY float64
	// Path is the outline a `<textPath>` writes along rather than a pen the
	// piece hangs on, and Offset is how far along it the writing starts. Both
	// are nil and zero for writing that has a pen instead, which is every
	// piece of writing a drawing has ever had.
	Path   *canvas.Path
	Offset float64
}
