package template

import "github.com/gabrielluizsf/antui/canvas"

// The measurements every widget is built out of, in units of the window's
// scale — the same `u` every painter below multiplies them by. They live here
// rather than inline because a widget's box is agreed on by at least three
// places: what it paints, what it measures, and what the interaction layer
// hit-tests. A number written twice is a number that will drift.
//
// Sizes are the whole widget; insets are gaps inside or around one.
const (
	// boxSize is the square a checkbox or a radio is drawn in — the tick's box,
	// the radio's circle, and the height of the row it sits on.
	boxSize = 18

	// switchTrackW and switchTrackH are a switch's pill: twice as wide as it is tall, so
	// the thumb travels the width of a diameter across the inside.
	switchTrackW = 36
	switchTrackH = 20

	// labelGap is the space between a control and the text beside it, the same
	// for a checkbox, a radio and a switch so the three rows line up.
	labelGap = 8

	// caretW is the room a dropdown keeps for the arrow at its right edge.
	caretW = 16

	// radiusMax caps a corner radius however tall the widget is: past this
	// the rounding is a circle, and a button would otherwise be a lozenge
	// stretched by its own label.
	radiusMax = 8

	// focusRing is how far the focus outline stands off the widget, and how
	// thick it is.
	focusRing = 2
	// focusRingRadius is the outline's own corner radius, a step past the
	// widget's so the two curves read as one shape.
	focusRingRadius = 6
)

// The rest are measured in rows of the built-in face rather than in whole
// units, so they follow the font if it is ever another size. `var` and not
// `const` because [canvas.FontWidth] and [canvas.FontHeight] are.
var (
	// controlW is the width of a slider and of an empty input — wide because
	// they hold a value rather than a label.
	controlW = 24 * canvas.FontWidth
	// sliderHeight is a slider's track: one font row. A progress bar's is two
	// thirds of that, a text area's is the same one row.
	sliderHeight    = 3 * canvas.FontHeight
	progressHeight  = 2 * canvas.FontHeight
	textAreaMinCols = 32 * canvas.FontWidth
)
