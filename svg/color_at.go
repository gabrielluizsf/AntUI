package svg

import "github.com/gabrielluizsf/antui/canvas"

// colorAt is the colour of the gradient at one point of the drawing, in the same
// coordinates the shape itself was written in.
func (s sampler) colorAt(x, y float64, current canvas.Color) canvas.Color {
	if len(s.stops) == 0 {
		return 0
	}
	gx, gy := s.toGradient.Map(x, y)
	return s.stopColor(s.spreadAt(s.parameter(gx, gy)), current)
}
