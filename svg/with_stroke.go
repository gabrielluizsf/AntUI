package svg

// withStroke reads a `stroke`: a colour, a reference to something elsewhere in
// the drawing, or nothing at all, the same way [Style.withFill] reads a fill.
func (s *Style) withStroke(raw string) {
	switch {
	case isCurrentColor(raw):
		s.Stroke, s.HasStroke, s.StrokeCurrent = 0, true, true
		s.strokeRef, s.strokeFallback, s.hasStrokeFallback, s.strokeFallbackCurrent = "", 0, false, false
		s.strokeGradient, s.strokePattern = nil, nil
	default:
		if c, ok := parsePaint(raw); ok {
			s.Stroke, s.HasStroke, s.StrokeCurrent = c, true, false
			s.strokeRef, s.strokeFallback, s.hasStrokeFallback, s.strokeFallbackCurrent = "", 0, false, false
			s.strokeGradient, s.strokePattern = nil, nil
		} else if id, ok := paintRef(raw); ok {
			s.strokeRef, s.strokeGradient, s.strokePattern = id, nil, nil
			s.strokeFallback, s.strokeFallbackCurrent, s.hasStrokeFallback = paintFallback(raw)
			s.Stroke, s.HasStroke, s.StrokeCurrent = 0, false, false
		} else {
			s.Stroke, s.HasStroke, s.StrokeCurrent = 0, false, false
			s.strokeRef, s.strokeFallback, s.hasStrokeFallback, s.strokeFallbackCurrent = "", 0, false, false
			s.strokeGradient, s.strokePattern = nil, nil
		}
	}
}
