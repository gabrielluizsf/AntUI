package svg

// withFill reads a `fill`: a colour to paint with, a reference to something
// elsewhere in the drawing, or nothing at all. `none` says so, and so does a
// paint this cannot read at all, which is not a reason to paint the shape
// black.
func (s *Style) withFill(raw string) {
	switch {
	case isCurrentColor(raw):
		s.Fill, s.HasFill, s.FillCurrent = 0, true, true
		s.fillRef, s.fillFallback, s.hasFillFallback, s.fillFallbackCurrent = "", 0, false, false
		s.fillGradient, s.fillPattern = nil, nil
	default:
		if c, ok := parsePaint(raw); ok {
			s.Fill, s.HasFill, s.FillCurrent = c, true, false
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
			s.Fill, s.HasFill, s.FillCurrent = 0, false, false
		} else {
			s.Fill, s.HasFill, s.FillCurrent = 0, false, false
			s.fillRef, s.fillFallback, s.hasFillFallback, s.fillFallbackCurrent = "", 0, false, false
			s.fillGradient, s.fillPattern = nil, nil
		}
	}
}
