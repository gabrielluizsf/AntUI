package canvas

// Gradient is the image a template paints into a rectangle with
// [Canvas.FillGradient]. Angle is radians measured clockwise from the top,
// matching "to top" being 0°; CenterX and CenterY are fractions of the
// rectangle the radial or conic gradient spins around. An Angle of zero and a
// centre of 0.5, 0.5 are the defaults a gradient without them wants.
type Gradient struct {
	Kind    GradientKind
	Angle   float64
	CenterX float64
	CenterY float64
	Shape   GradientShape
	Size    GradientSize
	Stops   []GradientStop
}
