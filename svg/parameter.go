package svg

// parameter is where a point of the gradient's own space lands along its line,
// 0 at the start and 1 at the far end. Past either end it keeps going, because
// what happens out there is the spread method's business and not this one's.
func (s sampler) parameter(x, y float64) float64 {
	if !s.linear {
		return s.radialParameter(x, y)
	}
	span := s.d.X*s.d.X + s.d.Y*s.d.Y
	if span == 0 {
		// A gradient whose ends are in the same place paints its last stop
		// everywhere, which is what SVG says rather than leaving it blank.
		return 1
	}
	return ((x-s.p0.X)*s.d.X + (y-s.p0.Y)*s.d.Y) / span
}
