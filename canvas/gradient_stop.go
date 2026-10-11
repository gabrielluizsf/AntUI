package canvas

// GradientStop is one colour stop of a gradient: a colour and where along the
// gradient line it changes, running from 0 to 1. The offsets a caller leaves
// out are spread evenly by [Canvas.FillGradient].
type GradientStop struct {
	Offset float64
	Color  Color
}
