package svg

import "strings"

// withVectorEffect reads `vector-effect`: whether the stroke is taken through
// everything that scaled the shape — which is what SVG does by default — or
// drawn at the width it says, in the pixels the drawing comes out at whatever
// scaled it.
//
// The other three values the property has — `non-scaling-size`, `non-rotation`
// and `fixed-position` — are refused, and a refusal that is said beats a stroke
// quietly painted as though nobody had asked for anything. All three are
// measured against the host coordinate system, the page the drawing is placed
// in, and this package paints one drawing onto one canvas with no page around
// it: the viewBox fit is the nearest thing to a host transform and it is not
// one. No browser draws them either, so a file that writes one is written
// against a draft rather than against the SVG a drawing can count on.
func (s *Style) withVectorEffect(raw string, warn func(string, ...any)) {
	switch keyword(raw) {
	case "none":
		s.NonScalingStroke = false
	case "non-scaling-stroke":
		s.NonScalingStroke = true
	default:
		s.NonScalingStroke = false
		if strings.TrimSpace(raw) != "" {
			warn("the vector-effect %q is not one this package draws (only none and non-scaling-stroke), so the stroke is scaled the way SVG scales it", raw)
		}
	}
}
