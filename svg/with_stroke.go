package svg

// withStroke reads a `stroke`: a colour, the keyword that takes the stroke of
// the naming shape inside a marker, a reference to something elsewhere in the
// drawing, or nothing at all, the same way [Style.withFill] reads a fill.
func (s *Style) withStroke(raw string) {
	switch {
	case isCurrentColor(raw):
		s.Stroke, s.HasStroke, s.StrokeCurrent, s.StrokeContext = 0, true, true, false
		s.strokeRef, s.strokeFallback, s.hasStrokeFallback, s.strokeFallbackCurrent = "", 0, false, false
		s.strokeGradient, s.strokePattern = nil, nil
	case isPaintKeyword(raw, "context-stroke"):
		// As with `context-fill`: the colour is the stroke of the element that
		// named the marker this is painted inside, known only when it is drawn.
		s.Stroke, s.HasStroke, s.StrokeCurrent, s.StrokeContext = 0, true, false, true
		s.strokeRef, s.strokeFallback, s.hasStrokeFallback, s.strokeFallbackCurrent = "", 0, false, false
		s.strokeGradient, s.strokePattern = nil, nil
	default:
		if c, ok := parsePaint(raw); ok {
			s.Stroke, s.HasStroke, s.StrokeCurrent, s.StrokeContext = c, true, false, false
			s.strokeRef, s.strokeFallback, s.hasStrokeFallback, s.strokeFallbackCurrent = "", 0, false, false
			s.strokeGradient, s.strokePattern = nil, nil
		} else if id, ok := paintRef(raw); ok {
			s.strokeRef, s.strokeGradient, s.strokePattern = id, nil, nil
			s.strokeFallback, s.strokeFallbackCurrent, s.hasStrokeFallback = paintFallback(raw)
			s.Stroke, s.HasStroke, s.StrokeCurrent, s.StrokeContext = 0, false, false, false
		} else {
			s.Stroke, s.HasStroke, s.StrokeCurrent, s.StrokeContext = 0, false, false, false
			s.strokeRef, s.strokeFallback, s.hasStrokeFallback, s.strokeFallbackCurrent = "", 0, false, false
			s.strokeGradient, s.strokePattern = nil, nil
		}
	}
}
