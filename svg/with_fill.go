package svg

// withFill reads a `fill`: a colour to paint with, the keyword that takes the
// colour of the naming shape inside a marker, a reference to something else in
// the drawing, or nothing at all. `none` says so, and so does a paint this
// cannot read at all, which is not a reason to paint the shape black.
func (s *Style) withFill(raw string) {
	switch {
	case isCurrentColor(raw):
		s.Fill, s.HasFill, s.FillCurrent, s.FillContext = 0, true, true, false
		s.fillRef, s.fillFallback, s.hasFillFallback, s.fillFallbackCurrent = "", 0, false, false
		s.fillGradient, s.fillPattern = nil, nil
	case isPaintKeyword(raw, "context-fill"):
		// The colour is not here and not the drawing's own either: it is the
		// fill of the element that named the marker this is painted inside, so
		// it is kept as a keyword and resolved where the marker is drawn.
		s.Fill, s.HasFill, s.FillCurrent, s.FillContext = 0, true, false, true
		s.fillRef, s.fillFallback, s.hasFillFallback, s.fillFallbackCurrent = "", 0, false, false
		s.fillGradient, s.fillPattern = nil, nil
	default:
		if c, ok := parsePaint(raw); ok {
			s.Fill, s.HasFill, s.FillCurrent, s.FillContext = c, true, false, false
			s.fillRef, s.fillFallback, s.hasFillFallback, s.fillFallbackCurrent = "", 0, false, false
			s.fillGradient, s.fillPattern = nil, nil
		} else if id, ok := paintRef(raw); ok {
			// The reference is kept whole: only a caller holding the whole
			// drawing can say what `#name` is, and the colour after it is the
			// one painted where the reference turns out to be nothing. Neither
			// the fill nor its colour is decided here — see
			// [Image.resolvePaints].
			s.fillRef, s.fillGradient, s.fillPattern = id, nil, nil
			s.fillFallback, s.fillFallbackCurrent, s.hasFillFallback = paintFallback(raw)
			s.Fill, s.HasFill, s.FillCurrent, s.FillContext = 0, false, false, false
		} else {
			s.Fill, s.HasFill, s.FillCurrent, s.FillContext = 0, false, false, false
			s.fillRef, s.fillFallback, s.hasFillFallback, s.fillFallbackCurrent = "", 0, false, false
			s.fillGradient, s.fillPattern = nil, nil
		}
	}
}
