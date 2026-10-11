package svg

import "github.com/gabrielluizsf/antui/canvas"

// A gradient is a paint rather than a colour: `fill="url(#g)"` names one
// elsewhere in the drawing, and the shape is painted by asking it what colour it
// is at each point rather than being painted one flat colour.
//
// The geometry is kept in the coordinates the gradient was written in and is not
// turned into a matrix at parse time, because SVG's own model is worth keeping:
// each number is a length in either user units or fractions of the shape's own
// box, and the two differ by whether x is divided by the box's width before it is
// used. A gradient written in a unit square and stretched to the shape it fills is
// exactly what a circular gradient in fractions of the bounding box is.
type gradient struct {
	// radial says a radial gradient; the others are linear. The kind is fixed by
	// the tag, so it is never inherited from the gradient an href points at.
	radial bool
	// userSpace says the geometry below is in user units rather than in fractions
	// of the shape's own bounding box.
	userSpace bool
	// stops are the colour changes along the gradient line, in the order they
	// were written.
	stops []gradientStop
	// transform is gradientTransform: it moves the gradient itself without
	// moving the shape it paints.
	transform canvas.Matrix
	hasXform  bool
	// spread is what happens past either end of the line: clamped, mirrored or
	// repeated.
	spread spread
	// set is which attributes the element actually wrote, which is what an href
	// consults: everything absent here is taken from the gradient it points at.
	set map[string]bool
	// href is the id of another gradient this one takes what it does not say from.
	href string
	// pct holds the attributes that were written as percentages, kept as they
	// were written until the units are settled. Which units they are is only
	// known once an href has been followed, and against a percentage the two are
	// worlds apart: a fraction of the shape, or a fraction of the drawing.
	pct map[string]pctRead

	// The linear line, from x1,y1 to x2,y2.
	x1, y1, x2, y2 float64
	// The radial circle, centred at cx,cy with radius r, and the focus fx,fy the
	// gradient is brightest at.
	cx, cy, r float64
	fx, fy    float64
	// viewport is the width and height a percentage written in user units is a
	// fraction of, which is the drawing's own and nothing to do with the shape
	// the gradient paints.
	viewport [2]float64
}
