package svg

import "github.com/gabrielluizsf/antui/canvas"

// sampler is a gradient with everything that does not depend on the pixel worked
// out once: where its line or its circle is, the stops with the offsets they left
// out filled in, and its spread. A pixel then costs a matrix, a distance and a
// mix, which is what makes a gradient in a drawing that is painted every frame
// affordable at all.
type sampler struct {
	// toGradient is what turns a point of the drawing's own coordinates into a
	// point of the gradient's: gradientTransform undone, and the shape's box
	// opened out to the unit square when the geometry is written in fractions of
	// it. That is the order SVG asks for, and reversing it is what makes a
	// gradientTransform mean the same thing in both of the units.
	toGradient canvas.Matrix
	// stops are the colour changes with the offsets they left out filled in.
	stops []gradientStop
	// spread is what happens past either end of the line.
	spread spread
	// linear says a line running from p0 along d; the others are a circle.
	linear bool
	p0, d  point
	// The circle, and the focus it is brightest at, which is never outside it.
	focus point
	c     point
	r     float64
}
